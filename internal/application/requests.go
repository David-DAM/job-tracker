package application

import "github.com/google/uuid"

type CreateJobRequest struct {
	Company     string `json:"company" binding:"required,min=2"`
	Position    string `json:"position" binding:"required,min=2"`
	Description string `json:"description" binding:"required,min=2"`
	Salary      int    `json:"salary"`
	Remote      bool   `json:"remote"`
	Url         string `json:"url"`
}

type UpdateJobRequest struct {
	Id          uuid.UUID `json:"id" binding:"required,uuid4"`
	Company     string    `json:"company" binding:"required,min=2"`
	Position    string    `json:"position" binding:"required,min=2"`
	Description string    `json:"description" binding:"required,min=2"`
	Salary      int       `json:"salary"`
	Remote      bool      `json:"remote"`
	Url         string    `json:"url"`
}
