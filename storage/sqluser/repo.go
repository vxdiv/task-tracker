package sqluser

import (
	"database/sql"
	"errors"
	"time"

	. "github.com/go-xorm/builder"

	"github.com/vxdiv/task-tracker/model"
	"github.com/vxdiv/task-tracker/storage/sqldbw"
)

var (
	ErrUserNotFound      = errors.New("user is not found")
	ErrUserAlreadyExists = errors.New("user is already exist")
)

type Repo interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Filter() *Filter
	Count(*Filter) (int, error)
	One(*Filter) (*model.User, error)
	List(*Filter) ([]*model.User, error)
}

func New(db *sql.DB) *SqlRepo {
	return &SqlRepo{sqldbw.New(db, "users")}
}

type SqlRepo struct {
	dbw *sqldbw.Wrapper
}

var _ Repo = &SqlRepo{}

func (r *SqlRepo) Create(m *model.User) error {
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

	if sqldbw.IsDuplicateEntryError(err) {
		return ErrUserAlreadyExists
	}

	return err
}

func (r *SqlRepo) Update(m *model.User) error {
	m.UpdatedAt = time.Now().UTC()
	return r.dbw.Update(Eq{
		"name":          m.Name,
		"email":         m.Email,
		"status":        m.Status,
		"password_hash": m.PasswordHash,
		"updated_at":    m.UpdatedAt,
	}, Eq{"id": m.ID})
}

func (r *SqlRepo) Filter() *Filter {
	return &Filter{
		repo: r,
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

func scanTo(m *model.User) (dst []interface{}) {
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

func (r *SqlRepo) Count(fr *Filter) (int, error) {
	return r.dbw.Count(fr.builder)
}

func (r *SqlRepo) One(fr *Filter) (*model.User, error) {
	user := &model.User{}
	err := r.dbw.Row(fr.builder, func(s sqldbw.RowScanner) error {
		return s.Scan(scanTo(user)...)
	})
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrUserNotFound
	default:
		return nil, err
	}

	return user, nil
}

func (r *SqlRepo) List(fr *Filter) ([]*model.User, error) {
	var list []*model.User
	err := r.dbw.Rows(fr.builder, func(s sqldbw.RowScanner) error {
		user := &model.User{}
		list = append(list, user)
		return s.Scan(scanTo(user)...)
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}
