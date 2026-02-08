package infrastructure

import (
	"job-tracker/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JobHandler struct {
	service *application.JobService
}

func NewJobHandler(s *application.JobService) *JobHandler {
	return &JobHandler{service: s}
}

func (h *JobHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/jobs", h.GetJobs)
	r.GET("/jobs/:id", h.GetJob)
	r.POST("/jobs", h.CreateJob)
	r.PUT("/jobs", h.UpdateJob)
	r.DELETE("/jobs/:id", h.DeleteJob)
	r.GET("/jobs/status/:status", h.GetJobsByStatus)
}

func (h *JobHandler) GetJobs(c *gin.Context) {
	jobs, err := h.service.GetAllJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetJobsByStatus(c *gin.Context) {
	status := c.Param("status")
	jobs, err := h.service.GetJobsByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := h.service.GetJob(uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
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

func (h *JobHandler) UpdateJob(c *gin.Context) {
	var request application.UpdateJobRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	job, err := h.service.UpdateJob(&request)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteJob(uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Job deleted successfully"})
}
