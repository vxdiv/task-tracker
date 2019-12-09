package app

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user is not found")
	ErrUserAlreadyExists = errors.New("user is already exist")
)

type UserRepo interface {
	Create(user *User) error
	Update(user *User) error
	GetByID(id int64) (*User, error)

	Count(filter UserFilter) (int, error)
	One(filter UserFilter) (*User, error)
	List(filter UserFilter) ([]*User, error)
}

type UserFilter struct {
	Name       string
	CreatedAt  TimeFilter
	Pagination Pager
}

type UserSvc struct {
	userRepo UserRepo
}

func NewUserSvc(repo UserRepo) *UserSvc {
	return &UserSvc{userRepo: repo}
}

func (svc *UserSvc) GetUser(id int64) (*User, error) {
	return svc.userRepo.GetByID(id)
}

func (svc *UserSvc) List(f UserFilter) (totalCount int, users []*User, err error) {
	if totalCount, err = svc.userRepo.Count(f); err != nil {
		return
	}

	if users, err = svc.userRepo.List(f); err != nil {
		return
	}

	return totalCount, users, nil
}

func (svc *UserSvc) Create(email, name, password string) error {
	user := &User{
		Email:  email,
		Name:   name,
		Status: UserStatusActive,
	}

	if err := user.SetPassword(password); err != nil {
		return err
	}

	if err := svc.userRepo.Create(user); err != nil {
		return err
	}

	return nil
}

func (svc *UserSvc) Update(id int64, name, status string) error {
	user, err := svc.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	user.Name = name
	user.Status = UserStatus(status)
	if err = svc.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}
