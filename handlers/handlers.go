package handlers

import (
	"database/sql"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/vxdiv/task-tracker/storage"
)

var (
	log *logrus.Entry

	users storage.UserRepo
)

func Init(db *sql.DB, e *echo.Echo, l *logrus.Entry) {
	log = l.WithField("component", "handlers")

	users = storage.NewUsers(db)

	e.Validator = &CustomValidator{validator: validator.New()}

	registerRoutes(e)
}

func registerRoutes(e *echo.Echo) {
	e.POST("/users", CreateUser)
	e.PUT("/users", UpdateUser)
	e.GET("/users/:id", GetUser)
	e.GET("/users", ListUser)
}
