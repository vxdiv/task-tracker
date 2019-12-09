package httpapi

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/vxdiv/task-tracker/app"
)

type UserCRUDHandler struct {
	userSvc *app.UserSvc
}

func RegisterUserCRUDHandler(e *echo.Echo, svc *app.UserSvc) {
	h := &UserCRUDHandler{userSvc: svc}
	e.POST("/users", h.Create)
	e.PUT("/users", h.Update)
	e.GET("/users/:id", h.View)
	e.GET("/users", h.List)
}

type RequestCreateUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RequestUpdateUser struct {
	ID     int64  `query:"ID" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Status string `json:"status" validate:"required, oneof=active disabled"`
}

type RequestListUsers struct {
	Name       string `query:"name"`
	CreateTime *TimeFilter
	Pagination *PagerRequestQuery
}

type ResponseUser struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Status  string `json:"status"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
}

func (h *UserCRUDHandler) Create(c echo.Context) error {
	req := &RequestCreateUser{}
	if err := parseRequest(c, req); err != nil {
		return err
	}

	err := h.userSvc.Create(req.Email, req.Name, req.Password)
	switch err {
	case nil:
	case app.ErrUserAlreadyExists:
		return BadRequestError(err)
	default:
		return err
	}

	return c.JSON(http.StatusCreated, "ok")
}

func (h *UserCRUDHandler) Update(c echo.Context) error {
	req := &RequestUpdateUser{}
	if err := parseRequest(c, req); err != nil {
		return err
	}

	err := h.userSvc.Update(req.ID, req.Name, req.Status)
	switch err {
	case nil:
	case app.ErrUserNotFound:
		return NotFoundError(err)
	default:
		return InternalServerError(err)
	}

	return c.JSON(http.StatusOK, "ok")
}

func (h *UserCRUDHandler) View(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.userSvc.GetUser(int64(id))
	switch err {
	case nil:
	case app.ErrUserNotFound:
		return NotFoundError(err)
	default:
		return InternalServerError(err)
	}

	return c.JSON(http.StatusOK, convertUser(user))
}

func (h *UserCRUDHandler) List(c echo.Context) error {
	req := &RequestListUsers{}
	if err := parseRequest(c, req); err != nil {
		return err
	}

	total, users, err := h.userSvc.List(app.UserFilter{
		Name:       req.Name,
		CreatedAt:  req.CreateTime,
		Pagination: req.Pagination,
	})
	if err != nil {
		return err
	}

	items := make([]ResponseUser, 0, len(users))
	for _, user := range users {
		items = append(items, convertUser(user))
	}

	return c.JSON(http.StatusOK, ResponseItems(total, req.Pagination, items))
}

func convertUser(user *app.User) ResponseUser {
	return ResponseUser{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Status:  user.Status.String(),
		Created: user.CreatedAt.Unix(),
		Updated: user.UpdatedAt.Unix(),
	}
}
