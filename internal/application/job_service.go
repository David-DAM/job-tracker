package application

import (
	"context"
	"job-tracker/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JobService struct {
	repository domain.JobRepository
	log        *domain.Logger
}

func NewJobService(repository domain.JobRepository, log *domain.Logger) *JobService {
	return &JobService{
		repository: repository,
		log:        log,
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

func (s *JobService) CreateJob(request *CreateJobRequest, ctx context.Context) (*domain.Job, error) {
	s.log.Ctx(ctx).Info("creating job")
	job := domain.NewJob(request.Company, request.Position, request.Description, request.Salary, request.Remote)
	err := s.repository.CreateJob(job)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to create job", zap.Error(err))
		return nil, err
	}
	s.log.Ctx(ctx).Info("job created", zap.String("job_id", job.Id.String()))
	return job, nil
}

func (s *JobService) UpdateJob(request *UpdateJobRequest, ctx context.Context) (*domain.Job, error) {
	s.log.Ctx(ctx).Info("updating job", zap.String("job_id", request.Id.String()))
	job, err := s.repository.GetJobById(request.Id.String())
	if err != nil {
		s.log.Ctx(ctx).Error("failed to get job to update", zap.Error(err))
		return nil, domain.ErrJobNotFound
	}
	job.Update(request.Company, request.Position, request.Description, request.Salary, request.Remote)
	err = s.repository.UpdateJob(job)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to update job", zap.Error(err))
		return nil, err
	}
	s.log.Ctx(ctx).Info("job updated", zap.String("job_id", job.Id.String()))
	return job, nil
}

func (s *JobService) DeleteJob(id uuid.UUID, ctx context.Context) error {
	s.log.Ctx(ctx).Info("deleting job", zap.String("job_id", id.String()))
	err := s.repository.DeleteJob(id.String())
	if err != nil {
		s.log.Ctx(ctx).Error("failed to delete job", zap.Error(err))
		return err
	}
	s.log.Ctx(ctx).Info("job deleted", zap.String("job_id", id.String()))
	return nil
}

func (s *JobService) GetAllJobs(ctx context.Context) ([]*domain.Job, error) {
	jobs, err := s.repository.GetAll()
	if err != nil {
		s.log.Ctx(ctx).Error("failed to get all jobs", zap.Error(err))
		return nil, err
	}
	return jobs, nil
}

func (s *JobService) GetJob(id uuid.UUID, ctx context.Context) (*domain.Job, error) {
	job, err := s.repository.GetJobById(id.String())
	if err != nil {
		s.log.Ctx(ctx).Error("failed to get job", zap.Error(err))
		return nil, domain.ErrJobNotFound
	}
	return job, nil
}

func (s *JobService) GetJobsByStatus(status string, ctx context.Context) ([]*domain.Job, error) {
	jobs, err := s.repository.GetJobsByStatus(status)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to get jobs by status", zap.Error(err))
		return nil, err
	}
	return jobs, nil
}
