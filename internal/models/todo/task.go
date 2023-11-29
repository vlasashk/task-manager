package todo

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
	Status      bool      `json:"status" validate:"required"`
}

func New(req Task) Task {
	req.ID = uuid.New().String()
	return req
}
