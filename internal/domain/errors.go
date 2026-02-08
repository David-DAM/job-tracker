package domain

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}

var ErrJobNotFound = errors.New("job not found")
var ErrJobAlreadyExists = errors.New("job already exists")
var ErrInvalidRequest = errors.New("invalid request")
var ErrInternalServer = errors.New("internal server error")
