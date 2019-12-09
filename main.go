package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/vxdiv/task-tracker/httpapi"
)

func main() {
	dbConn := mysqlConn(MysqlConfig{
		Host:     "localhost",
		Port:     3306,
		DBName:   "project",
		User:     "root",
		Password: "root",
	})
	defer func() {
		_ = dbConn.Close()
	}()

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	log := logrus.WithField("app", "task-tracker")

	httpApiServer := httpapi.NewServer(dbConn, log)
	if err := httpApiServer.Run(); err != nil {
		log.Fatal(err)
	}
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func mysqlConn(c MysqlConfig) *sql.DB {
	conn, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			c.User,
			c.Password,
			c.Host,
			c.Port,
			c.DBName,
		),
	)
	if err != nil {
		panic(err)
	}

	if err := conn.Ping(); err != nil {
		panic(err)
	}

	return conn
}
