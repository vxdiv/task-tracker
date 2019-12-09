package sqldb

import (
	"database/sql"
	"time"

	. "github.com/go-xorm/builder"

	"github.com/vxdiv/task-tracker/app"
)

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{NewDBWrapper(db, "users")}
}

type UserRepo struct {
	dbw *Wrapper
}

var _ app.UserRepo = &UserRepo{}

func (r *UserRepo) Create(m *app.User) error {
	var err error
	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = m.CreatedAt
	m.ID, err = r.dbw.Insert(Eq{
		"name":          m.Name,
		"email":         m.Email,
		"status":        m.Status,
		"password_hash": m.PasswordHash,
		"created_at":    m.CreatedAt,
		"updated_at":    m.UpdatedAt,
	})

	if IsDuplicateEntryError(err) {
		return app.ErrUserAlreadyExists
	}

	return err
}

func (r *UserRepo) Update(m *app.User) error {
	m.UpdatedAt = time.Now().UTC()
	return r.dbw.Update(Eq{
		"name":          m.Name,
		"email":         m.Email,
		"status":        m.Status,
		"password_hash": m.PasswordHash,
		"updated_at":    m.UpdatedAt,
	}, Eq{"id": m.ID})
}

func (r *UserRepo) selectBuilder() *Builder {
	return r.dbw.SelectBuilder(
		"id",
		"name",
		"email",
		"status",
		"password_hash",
		"created_at",
		"updated_at",
	)
}

func (r *UserRepo) scanTo(m *app.User) (dst []interface{}) {
	return []interface{}{
		&m.ID,
		&m.Name,
		&m.Email,
		&m.Status,
		&m.PasswordHash,
		&m.CreatedAt,
		&m.UpdatedAt,
	}
}

func (r *UserRepo) GetByID(id int64) (*app.User, error) {
	b := r.selectBuilder().Where(Eq{"id": id})

	return r.one(b)
}

func (r *UserRepo) filter(b *Builder, filter app.UserFilter) {
	if len(filter.Name) > 0 {
		b.Where(Eq{"name": filter.Name})
	}

	if filter.CreatedAt != nil {
		filterTime(b, filter.CreatedAt, "created_at")
	}
}

func (r *UserRepo) Count(f app.UserFilter) (int, error) {
	b := r.selectBuilder()
	r.filter(b, f)

	return r.dbw.Count(b)
}

func (r *UserRepo) One(f app.UserFilter) (*app.User, error) {
	b := r.selectBuilder()
	r.filter(b, f)

	return r.one(b)
}

func (r *UserRepo) one(b *Builder) (*app.User, error) {
	user := &app.User{}
	err := r.dbw.Row(b, func(s RowScanner) error {
		return s.Scan(r.scanTo(user)...)
	})
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, app.ErrUserNotFound
	default:
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) List(f app.UserFilter) ([]*app.User, error) {
	b := r.selectBuilder()
	r.filter(b, f)
	pagination(b, f.Pagination)

	var list []*app.User
	err := r.dbw.Rows(b, func(s RowScanner) error {
		user := &app.User{}
		list = append(list, user)
		return s.Scan(r.scanTo(user)...)
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}
