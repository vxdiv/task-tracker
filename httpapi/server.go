package httpapi

import (
	"database/sql"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"github.com/vxdiv/task-tracker/app"
	"github.com/vxdiv/task-tracker/sqldb"
)

type Server struct {
	e *echo.Echo
}

func NewServer(db *sql.DB, l *logrus.Entry) *Server {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	userRepo := sqldb.NewUserRepo(db)
	userSvc := app.NewUserSvc(userRepo)

	RegisterUserCRUDHandler(e, userSvc)
	srv := &Server{e: e}

	return srv
}

func (s *Server) Run() error {
	return s.e.Start(":1323")
}
