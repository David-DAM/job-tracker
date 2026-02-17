package application

import (
	"context"
	"job-tracker/internal/domain"

	"github.com/google/uuid"
)

type JobService struct {
	repository domain.JobRepository
	log        domain.Logger
}

func NewJobService(repository domain.JobRepository, log domain.Logger) *JobService {
	return &JobService{
		repository: repository,
		log:        log,
	}
}

func (s *JobService) CreateJob(request *CreateJobRequest, ctx context.Context) (*domain.Job, error) {
	s.log.Info(ctx, "creating job")
	job := domain.NewJob(request.Company, request.Position, request.Description, request.Salary, request.Remote)
	err := s.repository.CreateJob(job)
	if err != nil {
		s.log.Error(ctx, "failed to create job", err)
		return nil, err
	}
	s.log.Info(ctx, "job created", domain.Field{Key: "job_id", Value: job.Id.String()})
	return job, nil
}

func (s *JobService) UpdateJob(request *UpdateJobRequest, ctx context.Context) (*domain.Job, error) {
	s.log.Info(ctx, "updating job", domain.Field{Key: "job_id", Value: request.Id.String()})
	job, err := s.repository.GetJobById(request.Id.String())
	if err != nil {
		s.log.Error(ctx, "failed to get job to update", err)
		return nil, domain.ErrJobNotFound
	}
	job.Update(request.Company, request.Position, request.Description, request.Salary, request.Remote)
	err = s.repository.UpdateJob(job)
	if err != nil {
		s.log.Error(ctx, "failed to update job", err)
		return nil, err
	}
	s.log.Info(ctx, "job updated", domain.Field{Key: "job_id", Value: job.Id.String()})
	return job, nil
}

func (s *JobService) DeleteJob(id uuid.UUID, ctx context.Context) error {
	s.log.Info(ctx, "deleting job", domain.Field{Key: "job_id", Value: id.String()})
	err := s.repository.DeleteJob(id.String())
	if err != nil {
		s.log.Error(ctx, "failed to delete job", err)
		return err
	}
	s.log.Info(ctx, "job deleted", domain.Field{Key: "job_id", Value: id.String()})
	return nil
}

func (s *JobService) GetAllJobs(ctx context.Context) ([]*domain.Job, error) {
	jobs, err := s.repository.GetAll()
	if err != nil {
		s.log.Error(ctx, "failed to get all jobs", err)
		return nil, err
	}
	return jobs, nil
}

func (s *JobService) GetJob(id uuid.UUID, ctx context.Context) (*domain.Job, error) {
	job, err := s.repository.GetJobById(id.String())
	if err != nil {
		s.log.Error(ctx, "failed to get job", err)
		return nil, domain.ErrJobNotFound
	}
	return job, nil
}

func (s *JobService) GetJobsByStatus(status string, ctx context.Context) ([]*domain.Job, error) {
	jobs, err := s.repository.GetJobsByStatus(status)
	if err != nil {
		s.log.Error(ctx, "failed to get jobs by status", err)
		return nil, err
	}
	return jobs, nil
}
