package application

import (
	"fmt"
	"job-tracker/internal/domain"
	"time"

	"github.com/google/uuid"
)

type JobService struct {
	repository domain.JobRepository
}

func NewJobService(repository domain.JobRepository) *JobService {
	return &JobService{
		repository: repository,
	}
}

type CreateJobRequest struct {
	Company     string `json:"company" binding:"required,min=2"`
	Position    string `json:"position" binding:"required,min=2"`
	Description string `json:"description" binding:"required,min=2"`
	Salary      int    `json:"salary"`
	Remote      bool   `json:"remote"`
}

type UpdateJobRequest struct {
	Id          uuid.UUID `json:"id" binding:"required,uuid4"`
	Company     string    `json:"company" binding:"required,min=2"`
	Position    string    `json:"position" binding:"required,min=2"`
	Description string    `json:"description" binding:"required,min=2"`
	Salary      int       `json:"salary"`
	Remote      bool      `json:"remote"`
}

func (s *JobService) CreateJob(request *CreateJobRequest) (*domain.Job, error) {
	job := domain.NewJob(request.Company, request.Position, request.Description, request.Salary, request.Remote)
	err := s.repository.CreateJob(job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (s *JobService) UpdateJob(request *UpdateJobRequest) (*domain.Job, error) {
	job, _ := s.repository.GetJobById(request.Id.String())
	if job == nil {
		return nil, fmt.Errorf("job not found")
	}
	job.Company = request.Company
	job.Position = request.Position
	job.Description = request.Description
	job.Salary = request.Salary
	job.Remote = request.Remote
	job.UpdatedAt = time.Now()
	err := s.repository.UpdateJob(job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (s *JobService) DeleteJob(id uuid.UUID) error {
	return s.repository.DeleteJob(id.String())
}

func (s *JobService) GetAllJobs() ([]*domain.Job, error) {
	jobs, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *JobService) GetJob(id uuid.UUID) (*domain.Job, error) {
	job, err := s.repository.GetJobById(id.String())
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (s *JobService) GetJobsByStatus(status string) ([]*domain.Job, error) {
	jobs, err := s.repository.GetJobsByStatus(status)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
