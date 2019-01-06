package sqluser

import (
	"database/sql"

	"github.com/go-xorm/builder"
	"github.com/vxdiv/task-tracker/model"
	"github.com/vxdiv/task-tracker/storage"
	"github.com/vxdiv/task-tracker/storage/sqldbw"
)

type finder struct {
	dbw     *sqldbw.Wrapper
	builder *builder.Builder
}

var _ storage.UserFinder = &finder{}

func (f *finder) ByID(id int64) storage.UserFinder {
	f.builder.Where(builder.Eq{"id": id})

	return f
}

func (f *finder) ByName(name string) storage.UserFinder {
	f.builder.Where(builder.Eq{"name": name})

	return f
}

func (f *finder) ByEmail(email string) storage.UserFinder {
	f.builder.Where(builder.Eq{"email": email})

	return f
}

func (f *finder) CreatedAt(filter storage.TimeFilter) storage.UserFinder {
	if filter != nil {
		storage.FilterTime(f.builder, filter, "created_at")
	}

	return f
}

func (f *finder) Status(status string) storage.UserFinder {
	if len(status) > 0 {
		f.builder.Where(builder.Eq{"status": status})
	}

	return f
}

type itemScanner struct {
	item *model.User
}

func (is *itemScanner) ScanItem(s sqldbw.Scanner) error {
	return s.Scan(
		&is.item.ID,
		&is.item.Name,
		&is.item.Email,
		&is.item.Status,
		&is.item.PasswordHash,
		&is.item.CreatedAt,
		&is.item.UpdatedAt,
	)
}

func (f *finder) One() (*model.User, error) {
	user := &model.User{}
	err := f.dbw.LoadOne(f.builder, &itemScanner{user})
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, storage.ErrUserNotFound
	default:
		return nil, err
	}

	return user, nil
}

type listScanner struct {
	items []*model.User
}

func (ls *listScanner) ScanItem(s sqldbw.Scanner) error {
	item := &itemScanner{&model.User{}}
	if err := item.ScanItem(s); err != nil {
		return err
	}

	ls.items = append(ls.items, item.item)

	return nil
}

func (f *finder) List(limit storage.LimitFilter) ([]*model.User, int, error) {
	listScanner := &listScanner{}
	totalCount, err := f.dbw.LoadList(f.builder, listScanner, limit)

	return listScanner.items, totalCount, err
}
