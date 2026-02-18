package infrastructure

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"job-tracker/internal/domain"

	"golang.org/x/net/html"
	"golang.org/x/sync/semaphore"
)

type JobScrapper struct {
	rp   domain.JobRepository
	log  domain.Logger
	lock chan struct{}
}

func NewJobScrapper(rp domain.JobRepository, log domain.Logger) *JobScrapper {
	return &JobScrapper{
		rp:   rp,
		log:  log,
		lock: make(chan struct{}, 1),
	}
}

func (s *JobScrapper) InitScrape(ctx context.Context) error {

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	sem := semaphore.NewWeighted(20)

	for range ticker.C {
		s.log.Info(ctx, "scraping jobs")
		select {
		case s.lock <- struct{}{}:
			go func() {
				defer func() { <-s.lock }()
				if err := s.scrape(sem, ctx); err != nil {
					s.log.Error(ctx, "error scraping jobs", err)
					return
				}
			}()
		case <-ctx.Done():
			s.log.Info(ctx, "scrapper stopped")
			return nil
		default:
			continue
		}
		s.log.Info(ctx, "scraping finished")
	}
	return nil
}

func (s *JobScrapper) scrape(sem *semaphore.Weighted, ctx context.Context) error {

	jobs, err := s.rp.GetAll()
	if err != nil {
		s.log.Error(ctx, "error getting jobs", err)
		return err
	}

	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)
		if err := sem.Acquire(ctx, 1); err != nil {
			s.log.Error(ctx, "semaphore error", err)
			return err
		}

		go func(job *domain.Job) {

			defer func() { sem.Release(1); wg.Done() }()

			if err := s.run(job, ctx); err != nil {
				s.log.Error(ctx, "job failed", err)
			}
		}(job)
	}

	wg.Wait()
	return nil
}

func (s *JobScrapper) run(job *domain.Job, ctx context.Context) error {

	if job.Url == "" {
		return nil
	}

	resp, err := http.Get(job.Url)
	if err != nil {
		s.log.Error(ctx, "error making HTTP request", err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.log.Error(ctx, "error closing response body", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		s.log.Error(ctx, "non-OK HTTP status", nil)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Error(ctx, "error reading response body", err)
		return err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		s.log.Error(ctx, "error parsing HTML", err)
		return err
	}

	extractedInfo := s.extractJobInfo(doc)
	extractedStatus := s.extractJobStatus(doc)

	if extractedStatus == "" && extractedInfo == "" {
		return nil
	}

	if extractedInfo != "" {
		job.Description = extractedInfo
	}

	if extractedStatus != "" {
		job.Status = domain.JobStatusFromString(extractedStatus)
	}

	err = s.rp.UpdateJob(job)
	if err != nil {
		s.log.Error(ctx, "error updating job", err)
		return err
	}
	return nil
}

func (s *JobScrapper) extractJobInfo(n *html.Node) string {

	descriptionNode := s.findDescriptionNode(n)

	if descriptionNode == nil {
		return ""
	}

	var result strings.Builder
	s.traverse(&result, descriptionNode)

	out := result.String()
	out = strings.ReplaceAll(out, "\n ", "\n")
	out = strings.ReplaceAll(out, "  ", " ")

	return strings.TrimSpace(out)
}

func (s *JobScrapper) findDescriptionNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "description__text--rich") {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if res := s.findDescriptionNode(c); res != nil {
			return res
		}
	}

	return nil
}

func (s *JobScrapper) traverse(result *strings.Builder, n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "script", "style", "svg":
			return
		}

		if n.Data == "br" || n.Data == "p" || n.Data == "li" {
			result.WriteString("\n")
		}
	}

	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			result.WriteString(text)
			result.WriteString(" ")
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.traverse(result, c)
	}
}

func (s *JobScrapper) extractJobStatus(n *html.Node) string {
	var status string

	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {

			for _, attr := range n.Attr {
				if strings.Contains(strings.ToLower(attr.Key), "status") ||
					strings.Contains(strings.ToLower(attr.Val), "status") {
					if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
						status = strings.TrimSpace(n.FirstChild.Data)
						return
					}
				}
			}
		}
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)

			lowerText := strings.ToLower(text)
			if strings.Contains(lowerText, "open") ||
				strings.Contains(lowerText, "closed") ||
				strings.Contains(lowerText, "applied") ||
				strings.Contains(lowerText, "pending") ||
				strings.Contains(lowerText, "rejected") {
				if status == "" {
					status = text
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if status != "" {
				break
			}
			traverse(c)
		}
	}

	traverse(n)
	return strings.ToUpper(strings.TrimSpace(status))
}
