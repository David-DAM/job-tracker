package infrastructure

import (
	"job-tracker/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	service *application.JobService
}

func NewJobHandler(s *application.JobService) *JobHandler {
	return &JobHandler{service: s}
}

func (h *JobHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/jobs", h.GetJobs)
	r.POST("/jobs", h.CreateJob)
}

func (h *JobHandler) GetJobs(c *gin.Context) {
	jobs, err := h.service.GetAllJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) CreateJob(c *gin.Context) {

	var request application.CreateJobRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.service.CreateJob(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, job)
}
