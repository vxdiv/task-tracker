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

type userFinder struct {
	dbw     *sqldbw.Wrapper
	builder *builder.Builder
}

var _ storage.UserFinder = &userFinder{}

func (r Repo) Find() storage.UserFinder {
	return &userFinder{
		dbw:     r.dbw,
		builder: r.dbw.SelectBuilder("id", "name", "email", "status", "password_hash", "created_at", "updated_at")}
}

func (uf *userFinder) ByID(id int64) storage.UserFinder {
	uf.builder.Where(builder.Eq{"id": id})

	return uf
}

func (uf *userFinder) ByName(name string) storage.UserFinder {
	uf.builder.Where(builder.Eq{"name": name})

	return uf
}

func (uf *userFinder) ByEmail(email string) storage.UserFinder {
	uf.builder.Where(builder.Eq{"email": email})

	return uf
}

func (uf *userFinder) CreatedAt(filter storage.TimeFilter) storage.UserFinder {
	if filter != nil {
		storage.FilterTime(uf.builder, filter, "created_at")
	}

	return uf
}

func (uf *userFinder) Status(status string) storage.UserFinder {
	if len(status) > 0 {
		uf.builder.Where(builder.Eq{"status": status})
	}

	return uf
}

type userScan struct {
	user *model.User
}

func (us *userScan) ScanItem(s sqldbw.Scanner) error {
	return s.Scan(
		&us.user.ID,
		&us.user.Name,
		&us.user.Email,
		&us.user.Status,
		&us.user.PasswordHash,
		&us.user.CreatedAt,
		&us.user.UpdatedAt,
	)
}

func (uf *userFinder) One() (*model.User, error) {
	user := &model.User{}
	err := uf.dbw.LoadOne(uf.builder, &userScan{user})
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, storage.ErrUserNotFound
	default:
		return nil, err
	}

	return user, nil
}

type userListScan struct {
	items []*model.User
}

func (ul *userListScan) ScanItem(s sqldbw.Scanner) error {
	user := &model.User{}
	item := &userScan{user}
	if err := item.ScanItem(s); err != nil {
		return err
	}

	ul.items = append(ul.items, item.user)

	return nil
}

func (uf *userFinder) List(limit storage.LimitFilter) (list []*model.User, totalCount int, err error) {
	users := &userListScan{}
	totalCount, err = uf.dbw.LoadList(uf.builder, users, limit)

	return users.items, totalCount, err
}
