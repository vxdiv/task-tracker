package sqluser

import (
	"database/sql"
	"time"

	"github.com/go-xorm/builder"

	"github.com/vxdiv/task-tracker/model"
	"github.com/vxdiv/task-tracker/storage"
	"github.com/vxdiv/task-tracker/storage/sqldbw"
)

type Repo struct {
	dbw *sqldbw.Wrapper
}

var _ storage.UserRepo = &Repo{}

func New(db *sql.DB) *Repo {
	return &Repo{dbw: sqldbw.New(db, "users")}
}

func (r Repo) Create(user *model.User) error {
	var err error
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	user.ID, err = r.dbw.Insert(builder.Eq{
		"name":          user.Name,
		"email":         user.Email,
		"status":        user.Status,
		"password_hash": user.PasswordHash,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	})

	if sqldbw.IsDuplicateEntryError(err) {
		return storage.ErrUserAlreadyExists
	}

	return err
}

func (r Repo) Update(user *model.User) error {
	user.UpdatedAt = time.Now().UTC()
	return r.dbw.Update(builder.Eq{
		"name":          user.Name,
		"email":         user.Email,
		"status":        user.Status,
		"password_hash": user.PasswordHash,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}, builder.Eq{"id": user.ID})
}

func (r Repo) Find() storage.UserFinder {
	return &userFinder{
		dbw: r.dbw,
		builder: r.dbw.SelectBuilder(
			"id",
			"name",
			"email",
			"status",
			"password_hash",
			"created_at",
			"updated_at",
		)}
}
