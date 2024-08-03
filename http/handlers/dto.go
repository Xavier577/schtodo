package handlers

import "github.com/Xavier577/schtodo/pkg/date"

type IDParam struct {
	ID string `uri:"id" binding:"required"`
}

type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupPayload struct {
	LoginPayload
}

type CreateTodoPayload struct {
	Task     string        `json:"task" binding:"required"`
	IsTimed  bool          `json:"is_timed" default:"false"`
	Deadline date.DateTime `json:"deadline"`
}
