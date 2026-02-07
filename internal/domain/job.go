package domain

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	StatusPending   JobStatus = "PENDING"
	StatusApplied   JobStatus = "APPLIED"
	StatusInterview JobStatus = "INTERVIEW"
	StatusRejected  JobStatus = "REJECTED"
	StatusOffer     JobStatus = "OFFER"
)

type Job struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Company     string    `json:"company" validate:"required,min=2"`
	Position    string    `json:"position" validate:"required,min=2"`
	Description string    `json:"description"`
	Status      JobStatus `json:"status" validate:"required"`
	Salary      int       `json:"salary"`
	Remote      bool      `json:"remote"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type JobRepository interface {
	CreateJob(job *Job) error
	GetJobById(id string) (*Job, error)
	GetAll() ([]*Job, error)
	UpdateJob(job *Job) error
	DeleteJob(id string) error
	GetJobsByStatus(status string) ([]*Job, error)
}

func NewJob(company string, position string, description string, salary int, remote bool) *Job {
	return &Job{
		Id:          uuid.New(),
		Company:     company,
		Position:    position,
		Description: description,
		Status:      StatusPending,
		Salary:      salary,
		Remote:      remote,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
