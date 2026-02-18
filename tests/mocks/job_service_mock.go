package mocks

import (
	"job-tracker/internal/domain"

	"github.com/stretchr/testify/mock"
)

type JobRepositoryMock struct {
	mock.Mock
}

func (m *JobRepositoryMock) CreateJob(job *domain.Job) error {
	args := m.Called(job)
	return args.Error(0)
}

func (m *JobRepositoryMock) GetJobById(id string) (*domain.Job, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Job), args.Error(1)
}

func (m *JobRepositoryMock) GetAll() ([]*domain.Job, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Job), args.Error(1)
}

func (m *JobRepositoryMock) UpdateJob(job *domain.Job) error {
	args := m.Called(job)
	return args.Error(0)
}

func (m *JobRepositoryMock) DeleteJob(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *JobRepositoryMock) GetJobsByStatus(status domain.JobStatus) ([]*domain.Job, error) {
	args := m.Called(status)
	return args.Get(0).([]*domain.Job), args.Error(1)
}
