package application

import (
	"job-tracker/internal/domain"
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

func (s *JobService) CreateJob(request *CreateJobRequest) (*domain.Job, error) {
	job := domain.NewJob(request.Company, request.Position, request.Description, request.Salary, request.Remote)
	err := s.repository.CreateJob(job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (s *JobService) GetAllJobs() ([]*domain.Job, error) {
	jobs, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
