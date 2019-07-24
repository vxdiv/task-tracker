package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/vxdiv/task-tracker/model"
	"github.com/vxdiv/task-tracker/storage/sqluser"
)

type RequestCreateUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func CreateUser(c echo.Context) error {
	req := &RequestCreateUser{}
	if err := parseRequest(c, req); err != nil {
		return err
	}

	user := &model.User{
		Email:  req.Email,
		Name:   req.Name,
		Status: model.UserStatusActive,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return InternalServerError(err)
	}

	err := users.Create(user)
	switch err {
	case nil:
	case sqluser.ErrUserAlreadyExists:
		return BadRequestError(err)
	default:
		return InternalServerError(err)
	}

	return c.JSON(http.StatusCreated, "ok")
}

type RequestUpdateUser struct {
	ID     int64  `query:"ID" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Status string `json:"status" validate:"required, oneof=active disabled"`
}

func UpdateUser(c echo.Context) error {
	req := &RequestUpdateUser{}
	if err := parseRequest(c, req); err != nil {
		return err
	}

	user, err := users.Filter().ID(req.ID).One()
	switch err {
	case nil:
	case sqluser.ErrUserNotFound:
		return NotFoundError(err)
	default:
		return InternalServerError(err)
	}

	user.Name = req.Name
	user.Status = model.UserStatus(req.Status)
	if err := users.Update(user); err != nil {
		return InternalServerError(err)
	}

	return c.JSON(http.StatusOK, "ok")
}

type ResponseUser struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Status  string `json:"status"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
}

func GetUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := users.One(users.Filter().ID(int64(id)))
	switch err {
	case nil:
	case sqluser.ErrUserNotFound:
		return NotFoundError(err)
	default:
		return InternalServerError(err)
	}

	return c.JSON(http.StatusOK, convertUserFromModel(user))
}

type UserFilter struct {
	pager      *Pager
	timeFilter *TimeFilter
}

func (uf *UserFilter) Parse(c echo.Context) (err error) {
	if uf.pager, err = BindPager(c); err != nil {
		return err
	}

	if uf.timeFilter, err = BindTimeFilter(c); err != nil {
		return err
	}

	return nil
}

func ListUser(c echo.Context) error {
	filter := UserFilter{}
	if err := filter.Parse(c); err != nil {
		return BadRequestError(err)
	}

	finder := users.Filter().CreatedAt(filter.timeFilter)
	total, err := users.Count(finder)
	if err != nil {
		return InternalServerError(err)
	}

	users, err := users.List(finder.Limit(filter.pager))
	if err != nil {
		return InternalServerError(err)
	}

	items := make([]ResponseUser, 0, len(users))
	for _, user := range users {
		items = append(items, convertUserFromModel(user))
	}

	return c.JSON(http.StatusOK, ResponseItems(total, filter.pager, items))
}

func convertUserFromModel(user *model.User) ResponseUser {
	return ResponseUser{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Status:  user.Status.String(),
		Created: user.CreatedAt.Unix(),
		Updated: user.UpdatedAt.Unix(),
	}
}
