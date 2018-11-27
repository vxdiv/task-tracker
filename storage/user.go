package storage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-xorm/builder"

	"github.com/vxdiv/task-tracker/model"
)

var (
	ErrUserNotFound      = errors.New("user is not found")
	ErrUserAlreadyExists = errors.New("user is already exist")
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
	All(limit LimitFilter) ([]*model.User, int, error)
}

const tableUsers = "users"

type userRepo struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) UserRepo {
	return &userRepo{db}
}

func (repo userRepo) Create(user *model.User) error {
	var err error
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	user.ID, err = insert(repo.db, builder.Eq{
		"name":          user.Name,
		"email":         user.Email,
		"status":        user.Status,
		"password_hash": user.PasswordHash,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}, tableUsers)

	if isDuplicateEntryError(err) {
		return ErrUserAlreadyExists
	}

	return err
}

func (repo userRepo) Update(user *model.User) error {
	user.UpdatedAt = time.Now().UTC()
	return update(repo.db, builder.Eq{
		"name":          user.Name,
		"email":         user.Email,
		"status":        user.Status,
		"password_hash": user.PasswordHash,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}, builder.Eq{"id": user.ID}, tableUsers)
}

func (repo userRepo) Find() UserFinder {
	return &userFinder{builder: builder.Dialect(dbDialect).Select(
		"id",
		"name",
		"email",
		"status",
		"password_hash",
		"created_at",
		"updated_at",
	).From(tableUsers), db: repo.db}
}

type userFinder struct {
	builder *builder.Builder
	db      *sql.DB
}

func (finder *userFinder) ByID(id int64) UserFinder {
	finder.builder.Where(builder.Eq{"id": id})

	return finder
}

func (finder *userFinder) ByName(name string) UserFinder {
	finder.builder.Where(builder.Eq{"name": name})

	return finder
}

func (finder *userFinder) ByEmail(email string) UserFinder {
	finder.builder.Where(builder.Eq{"email": email})

	return finder
}

func (finder *userFinder) CreatedAt(filter TimeFilter) UserFinder {
	if filter != nil {
		filterTime(finder.builder, filter, "created_at")
	}

	return finder
}

func (finder *userFinder) Status(status string) UserFinder {
	if len(status) > 0 {
		finder.builder.Where(builder.Eq{"status": status})
	}

	return finder
}

type userScan struct {
	User *model.User
}

func (us *userScan) ScanItem(s Scanner) error {
	return s.Scan(
		&us.User.ID,
		&us.User.Name,
		&us.User.Email,
		&us.User.Status,
		&us.User.PasswordHash,
		&us.User.CreatedAt,
		&us.User.UpdatedAt,
	)
}

func (finder *userFinder) One() (*model.User, error) {
	user := &model.User{}
	err := loadOne(finder.db, finder.builder, &userScan{user})
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrUserNotFound
	default:
		return nil, err
	}

	return user, nil

}

type userListScan struct {
	items []*model.User
}

func (ul *userListScan) ScanItem(s Scanner) error {
	user := &model.User{}
	item := &userScan{user}
	if err := item.ScanItem(s); err != nil {
		return err
	}

	ul.items = append(ul.items, item.User)

	return nil
}

func (finder *userFinder) All(limit LimitFilter) ([]*model.User, int, error) {
	users := &userListScan{}
	totalCount, err := loadAll(finder.db, finder.builder, users, limit)

	return users.items, totalCount, err
}
