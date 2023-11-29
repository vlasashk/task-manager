package todo

import (
	"github.com/google/uuid"
)

type Task struct {
	ID string `json:"id,omitempty" validate:"required"`
	TaskReq
}

type TaskReq struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	DueDate     string `json:"due_date" validate:"required"`
	Status      bool   `json:"status" validate:"required"`
}

func New(req TaskReq) Task {
	return Task{
		ID:      uuid.New().String(),
		TaskReq: req,
	}
}
