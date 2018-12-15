package storage

import (
	"errors"

	"github.com/vxdiv/task-tracker/model"
)

var (
	ErrUserNotFound      = errors.New("users is not found")
	ErrUserAlreadyExists = errors.New("users is already exist")
)

type UserRepo interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Find() UserFinder
}

type UserFinder interface {
	ByID(id int64) UserFinder
	ByEmail(email string) UserFinder
	ByName(name string) UserFinder

	CreatedAt(filter TimeFilter) UserFinder
	Status(status string) UserFinder

	One() (*model.User, error)
	List(limit LimitFilter) (list []*model.User, totalCount int, err error)
}
