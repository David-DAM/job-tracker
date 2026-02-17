package infrastructure

import (
	"job-tracker/internal/domain"
	"job-tracker/internal/infrastructure"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&domain.Job{})
	assert.NoError(t, err)

	return db
}

func TestCreateAndGetJob(t *testing.T) {
	db := setupTestDB(t)
	repo := infrastructure.NewJobRepository(db)

	job := domain.NewJob("Google", "Backend", "Go dev", 100000, true)

	err := repo.CreateJob(job)
	assert.NoError(t, err)
	assert.NotEmpty(t, job.Id)

	found, err := repo.GetJobById(job.Id.String())
	assert.NoError(t, err)
	assert.Equal(t, "Google", found.Company)
	assert.Equal(t, "Backend", found.Position)
}

func TestGetAllJobs(t *testing.T) {
	db := setupTestDB(t)
	repo := infrastructure.NewJobRepository(db)

	job1 := domain.NewJob("Google", "Backend", "Go", 100, true)
	job2 := domain.NewJob("Amazon", "Java", "Spring", 200, false)

	_ = repo.CreateJob(job1)
	_ = repo.CreateJob(job2)

	jobs, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, jobs, 2)
}

func TestUpdateJob(t *testing.T) {
	db := setupTestDB(t)
	repo := infrastructure.NewJobRepository(db)

	job := domain.NewJob("Google", "Backend", "Go", 100, true)
	_ = repo.CreateJob(job)

	job.Company = "Meta"
	err := repo.UpdateJob(job)

	assert.NoError(t, err)

	updated, _ := repo.GetJobById(job.Id.String())
	assert.Equal(t, "Meta", updated.Company)
}

func TestDeleteJob(t *testing.T) {
	db := setupTestDB(t)
	repo := infrastructure.NewJobRepository(db)

	job := domain.NewJob("Google", "Backend", "Go", 100, true)
	_ = repo.CreateJob(job)

	err := repo.DeleteJob(job.Id.String())
	assert.NoError(t, err)

	_, err = repo.GetJobById(job.Id.String())
	assert.Error(t, err)
}

func TestDeleteJob_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := infrastructure.NewJobRepository(db)

	err := repo.DeleteJob("non-existing")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrJobNotFound, err)
}
