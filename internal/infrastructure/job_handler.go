package infrastructure

import (
	"errors"
	"job-tracker/internal/application"
	"job-tracker/internal/domain"
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
	isError := hasError(err, c)
	if isError {
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetJobsByStatus(c *gin.Context) {
	status := c.Param("status")
	jobs, err := h.service.GetJobsByStatus(status)
	isError := hasError(err, c)
	if isError {
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id, ok := parseUUID(c, c.Param("id"))

	if !ok {
		return
	}
	job, err := h.service.GetJob(id)
	isError := hasError(err, c)
	if isError {
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) CreateJob(c *gin.Context) {

	var request application.CreateJobRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: domain.ErrInvalidRequest.Error()})
		return
	}

	job, err := h.service.CreateJob(&request)
	isError := hasError(err, c)
	if isError {
		return
	}
	c.JSON(http.StatusCreated, job)
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	var request application.UpdateJobRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: domain.ErrInvalidRequest.Error()})
		return
	}
	job, err := h.service.UpdateJob(&request)
	isError := hasError(err, c)
	if isError {
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) DeleteJob(c *gin.Context) {

	id, ok := parseUUID(c, c.Param("id"))

	if !ok {
		return
	}
	err := h.service.DeleteJob(id)
	isError := hasError(err, c)
	if isError {
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Job deleted successfully"})
}

func parseUUID(c *gin.Context, idStr string) (uuid.UUID, bool) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: domain.ErrInvalidRequest.Error()})
		return uuid.Nil, false
	}
	return id, true
}

func hasError(err error, c *gin.Context) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, domain.ErrJobNotFound):
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Error: domain.ErrJobNotFound.Error()})
	case errors.Is(err, domain.ErrInvalidRequest):
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: domain.ErrInvalidRequest.Error()})
	default:
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: domain.ErrInternalServer.Error()})
	}
	return true
}
