package app

import "errors"

var (
	ErrTaskNotFound      = errors.New("task is not found")
	ErrTaskAlreadyExists = errors.New("task is already exist")
)

type TaskRepo interface {
	Create(task *Task) error
	Update(task *Task) error
	GetByID(id int64) (*Task, error)

	Count(filter TaskFilter) (int, error)
	One(filter TaskFilter) (*Task, error)
	List(filter TaskFilter) ([]*Task, error)
}

type TaskFilter struct {
	CreatedAt  TimeFilter
	Pagination Pager
}
