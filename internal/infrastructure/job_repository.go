package infrastructure

import (
	"job-tracker/internal/domain"

	"gorm.io/gorm"
)

type JobRepositoryImpl struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) domain.JobRepository {
	return &JobRepositoryImpl{
		db: db,
	}
}
func (r *JobRepositoryImpl) CreateJob(job *domain.Job) error {
	return r.db.Create(job).Error
}

func (r *JobRepositoryImpl) GetJobById(id string) (*domain.Job, error) {
	var job domain.Job
	err := r.db.First(&job, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepositoryImpl) GetAll() ([]*domain.Job, error) {
	var jobs []*domain.Job
	err := r.db.Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepositoryImpl) UpdateJob(job *domain.Job) error {
	return r.db.Save(job).Error
}

func (r *JobRepositoryImpl) DeleteJob(id string) error {
	result := r.db.Delete(&domain.Job{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrJobNotFound
	}
	return nil
}

func (r *JobRepositoryImpl) GetJobsByStatus(status domain.JobStatus) ([]*domain.Job, error) {
	var jobs []*domain.Job
	err := r.db.Where("status = ?", status).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
