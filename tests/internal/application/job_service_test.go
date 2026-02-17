package application

import (
	"context"
	"errors"
	"job-tracker/internal/application"
	"job-tracker/internal/domain"
	"job-tracker/tests/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func InitAppTest() (*mocks.JobRepositoryMock, *application.JobService) {
	var repo = new(mocks.JobRepositoryMock)
	var logger = &mocks.LoggerMock{}
	return repo, application.NewJobService(repo, logger)
}

func TestCreateJob(t *testing.T) {

	repo, service := InitAppTest()

	req := &application.CreateJobRequest{
		Company:     "Google",
		Position:    "Backend",
		Description: "Go dev",
		Salary:      100000,
		Remote:      true,
	}

	repo.On("CreateJob", mock.AnythingOfType("*domain.Job")).Return(nil)

	job, err := service.CreateJob(req, context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, job)
	assert.Equal(t, "Google", job.Company)

	repo.AssertExpectations(t)
}

func TestUpdateJob(t *testing.T) {

	repo, service := InitAppTest()

	jobID := uuid.New()
	req := &application.UpdateJobRequest{
		Id:          jobID,
		Company:     "Amazon",
		Position:    "Go Dev",
		Description: "Backend work",
		Salary:      120000,
		Remote:      true,
	}

	existingJob := domain.NewJob("OldCo", "OldPos", "OldDesc", 90000, false)
	existingJob.Id = jobID

	repo.On("GetJobById", jobID.String()).Return(existingJob, nil)
	repo.On("UpdateJob", existingJob).Return(nil)

	job, err := service.UpdateJob(req, context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "Amazon", job.Company)
	repo.AssertExpectations(t)
}

func TestUpdateJob_NotFound(t *testing.T) {

	repo, service := InitAppTest()

	jobID := uuid.New()
	req := &application.UpdateJobRequest{Id: jobID}

	repo.On("GetJobById", jobID.String()).Return(nil, domain.ErrJobNotFound)

	job, err := service.UpdateJob(req, context.Background())

	assert.ErrorIs(t, err, domain.ErrJobNotFound)
	assert.Nil(t, job)
	repo.AssertExpectations(t)
}

func TestDeleteJob(t *testing.T) {

	repo, service := InitAppTest()

	jobID := uuid.New()
	repo.On("DeleteJob", jobID.String()).Return(nil)

	err := service.DeleteJob(jobID, context.Background())
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteJob_NotFound(t *testing.T) {

	repo, service := InitAppTest()

	jobID := uuid.New()
	repo.On("DeleteJob", jobID.String()).Return(domain.ErrJobNotFound)

	err := service.DeleteJob(jobID, context.Background())
	assert.ErrorIs(t, err, domain.ErrJobNotFound)
	repo.AssertExpectations(t)
}

func TestGetJob(t *testing.T) {

	repo, service := InitAppTest()

	jobID := uuid.New()
	existingJob := domain.NewJob("Google", "Backend", "Go dev", 100000, true)
	existingJob.Id = jobID

	repo.On("GetJobById", jobID.String()).Return(existingJob, nil)

	job, err := service.GetJob(jobID, context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "Google", job.Company)
	repo.AssertExpectations(t)
}

func TestGetJob_NotFound(t *testing.T) {

	repo, service := InitAppTest()

	jobID := uuid.New()
	repo.On("GetJobById", jobID.String()).Return(nil, domain.ErrJobNotFound)

	job, err := service.GetJob(jobID, context.Background())
	assert.ErrorIs(t, err, domain.ErrJobNotFound)
	assert.Nil(t, job)
	repo.AssertExpectations(t)
}

func TestGetAllJobs(t *testing.T) {

	repo, service := InitAppTest()

	jobs := []*domain.Job{
		domain.NewJob("A", "X", "desc", 1, false),
		domain.NewJob("B", "Y", "desc", 2, true),
	}

	repo.On("GetAll").Return(jobs, nil)

	result, err := service.GetAllJobs(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestGetJobsByStatus(t *testing.T) {

	repo, service := InitAppTest()

	status := "open"
	jobs := []*domain.Job{
		domain.NewJob("A", "X", "desc", 1, false),
	}
	repo.On("GetJobsByStatus", status).Return(jobs, nil)

	result, err := service.GetJobsByStatus(status, context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestGetJobsByStatus_Error(t *testing.T) {

	repo, service := InitAppTest()
	status := "open"
	repo.On("GetJobsByStatus", status).Return([]*domain.Job{}, errors.New("db error"))

	result, err := service.GetJobsByStatus(status, context.Background())
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}
