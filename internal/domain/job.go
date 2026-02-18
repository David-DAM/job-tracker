package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusUnknown   JobStatus = "UNKNOWN"
	JobStatusOpen      JobStatus = "OPEN"
	JobStatusClosed    JobStatus = "CLOSED"
	JobStatusPending   JobStatus = "PENDING"
	JobStatusApplied   JobStatus = "APPLIED"
	JobStatusInterview JobStatus = "INTERVIEW"
	JobStatusRejected  JobStatus = "REJECTED"
	JobStatusOffer     JobStatus = "OFFER"
)

type Job struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Company     string    `json:"company" validate:"required,min=2"`
	Position    string    `json:"position" validate:"required,min=2"`
	Description string    `json:"description"`
	Status      JobStatus `json:"status" validate:"required"`
	Salary      int       `json:"salary"`
	Remote      bool      `json:"remote"`
	Url         string    `json:"url"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type JobRepository interface {
	CreateJob(job *Job) error
	GetJobById(id string) (*Job, error)
	GetAll() ([]*Job, error)
	UpdateJob(job *Job) error
	DeleteJob(id string) error
	GetJobsByStatus(status JobStatus) ([]*Job, error)
}

func NewJob(company string, position string, description string, salary int, remote bool, url string) *Job {
	return &Job{
		Id:          uuid.New(),
		Company:     company,
		Position:    position,
		Description: description,
		Status:      JobStatusPending,
		Salary:      salary,
		Remote:      remote,
		Url:         url,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func JobStatusFromString(status string) JobStatus {
	status = strings.ToUpper(status)
	switch status {
	case "OPEN":
		return JobStatusOpen
	case "CLOSED":
		return JobStatusClosed
	case "PENDING":
		return JobStatusPending
	case "APPLIED":
		return JobStatusApplied
	case "INTERVIEW":
		return JobStatusInterview
	default:
		return JobStatusUnknown
	}
}

func (j *Job) Update(company string, position string, description string, salary int, remote bool, url string) {
	j.UpdatedAt = time.Now()
	j.Company = company
	j.Position = position
	j.Description = description
	j.Salary = salary
	j.Remote = remote
	j.Url = url
}
