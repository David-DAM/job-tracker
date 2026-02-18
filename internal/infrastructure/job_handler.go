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
	logger  domain.Logger
}

func NewJobHandler(s *application.JobService, logger domain.Logger) *JobHandler {
	return &JobHandler{service: s, logger: logger}
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
	h.logger.Info(c.Request.Context(), "getting all jobs")
	jobs, err := h.service.GetAllJobs(c.Request.Context())
	isError := hasError(err, c)
	if isError {
		h.logger.Error(c.Request.Context(), "failed to get all jobs", err)
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetJobsByStatus(c *gin.Context) {
	h.logger.Info(c.Request.Context(), "getting jobs by status")
	status := c.Param("status")
	jobs, err := h.service.GetJobsByStatus(domain.JobStatusFromString(status), c.Request.Context())
	isError := hasError(err, c)
	if isError {
		h.logger.Error(c.Request.Context(), "failed to get jobs by status", err)
		return
	}
	h.logger.Info(c.Request.Context(), "jobs fetched successfully")
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	h.logger.Info(c.Request.Context(), "getting job by id")
	id, ok := parseUUID(c, c.Param("id"))
	if !ok {
		return
	}
	job, err := h.service.GetJob(id, c.Request.Context())
	isError := hasError(err, c)
	if isError {
		h.logger.Error(c.Request.Context(), "failed to get job", err)
		return
	}
	h.logger.Info(c.Request.Context(), "job fetched successfully")
	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	h.logger.Info(c.Request.Context(), "creating job")
	var request application.CreateJobRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error(c.Request.Context(), "invalid request body", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: domain.ErrInvalidRequest.Error()})
		return
	}

	job, err := h.service.CreateJob(&request, c.Request.Context())
	isError := hasError(err, c)
	if isError {
		h.logger.Error(c.Request.Context(), "failed to create job", err)
		return
	}
	h.logger.Info(c.Request.Context(), "job created successfully")
	c.JSON(http.StatusCreated, job)
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	h.logger.Info(c.Request.Context(), "updating job")
	var request application.UpdateJobRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error(c.Request.Context(), "invalid request body", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: domain.ErrInvalidRequest.Error()})
		return
	}
	job, err := h.service.UpdateJob(&request, c.Request.Context())
	isError := hasError(err, c)
	if isError {
		h.logger.Error(c.Request.Context(), "failed to update job", err)
		return
	}
	h.logger.Info(c.Request.Context(), "job updated successfully")
	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	h.logger.Info(c.Request.Context(), "deleting job")
	id, ok := parseUUID(c, c.Param("id"))

	if !ok {
		h.logger.Error(c.Request.Context(), "invalid job id provided", errors.New(c.Param("id")))
		return
	}
	err := h.service.DeleteJob(id, c.Request.Context())
	isError := hasError(err, c)
	if isError {
		h.logger.Error(c.Request.Context(), "failed to delete job", err)
		return
	}
	h.logger.Info(c.Request.Context(), "job deleted successfully")
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
