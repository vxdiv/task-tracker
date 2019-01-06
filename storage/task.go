package storage

import (
	"errors"

	"github.com/vxdiv/task-tracker/model"
)

var (
	ErrTaskNotFound      = errors.New("task is not found")
	ErrTaskAlreadyExists = errors.New("task is already exist")
)

type TaskRepo interface {
	Create(task *model.Task) error
	Update(task *model.Task) error
	Find() TaskFinder
}

type TaskFinder interface {
	ByID(id int64) TaskFinder

	CreatedAt(filter TimeFilter) TaskFinder
	Status(status string) TaskFinder
	Type(typ string) TaskFinder
	OwnerID(id int64) TaskFinder
	AssignedID(id int64) TaskFinder
	Resolution(resolution string) TaskFinder
	Priority(priority string) TaskFinder
	NameLike(name string) TaskFinder
	DueDate(filter TimeFilter) TaskFinder

	One() (*model.Task, error)
	List(limit LimitFilter) (list []*model.Task, totalCount int, err error)
}
