package tasks

import "errors"

var ErrNotFound = errors.New("task not found")

type Store interface {
	GetAll(limit, offset int) ([]Task, int, error) // Expects limit/offset, returns tasks, total_count, error
	GetByID(id int) (Task, error)
	Create(title string) (Task, error)
	Update(id int, title *string, done *bool) (Task, error)
	Delete(id int) error
}
