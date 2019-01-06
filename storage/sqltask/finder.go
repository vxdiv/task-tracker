package sqltask

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

var _ storage.TaskFinder = &finder{}

func (f *finder) ByID(id int64) storage.TaskFinder {
	f.builder.Where(builder.Eq{"id": id})

	return f
}

func (f *finder) CreatedAt(filter storage.TimeFilter) storage.TaskFinder {
	if filter != nil {
		storage.FilterTime(f.builder, filter, "created_at")
	}

	return f
}

func (f *finder) Status(status string) storage.TaskFinder {
	if len(status) > 0 {
		f.builder.Where(builder.Eq{"status": status})
	}

	return f
}

func (f *finder) Type(typ string) storage.TaskFinder {
	if len(typ) > 0 {
		f.builder.Where(builder.Eq{"type": typ})
	}

	return f
}

func (f *finder) OwnerID(id int64) storage.TaskFinder {
	if id > 0 {
		f.builder.Where(builder.Eq{"owner_id": id})
	}

	return f
}

func (f *finder) AssignedID(id int64) storage.TaskFinder {
	if id > 0 {
		f.builder.Where(builder.Eq{"assigned_id": id})
	}

	return f
}

func (f *finder) Resolution(resolution string) storage.TaskFinder {
	if len(resolution) > 0 {
		f.builder.Where(builder.Eq{"resolution": resolution})
	}

	return f
}

func (f *finder) Priority(priority string) storage.TaskFinder {
	if len(priority) > 0 {
		f.builder.Where(builder.Eq{"priority": priority})
	}

	return f
}

func (f *finder) NameLike(name string) storage.TaskFinder {
	if len(name) > 0 {
		f.builder.Where(builder.Like{"name", "%" + name})
	}

	return f
}

func (f *finder) DueDate(filter storage.TimeFilter) storage.TaskFinder {
	if filter != nil {
		storage.FilterTime(f.builder, filter, "due_date")
	}

	return f
}

type itemScanner struct {
	item *model.Task
}

func (ts *itemScanner) ScanItem(s sqldbw.Scanner) error {
	return s.Scan(
		&ts.item.ID,
		&ts.item.Name,
		&ts.item.Description,
		&ts.item.Type,
		&ts.item.OwnerID,
		&ts.item.AssignedID,
		&ts.item.Status,
		&ts.item.DueDate,
		&ts.item.Resolution,
		&ts.item.Priority,
		&ts.item.CreatedAt,
		&ts.item.UpdatedAt,
	)
}

func (f *finder) One() (*model.Task, error) {
	task := &model.Task{}
	err := f.dbw.LoadOne(f.builder, &itemScanner{task})
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, storage.ErrTaskNotFound
	default:
		return nil, err
	}

	return task, nil
}

type listScanner struct {
	items []*model.Task
}

func (tl *listScanner) ScanItem(s sqldbw.Scanner) error {
	task := &model.Task{}
	item := &itemScanner{task}
	if err := item.ScanItem(s); err != nil {
		return err
	}

	tl.items = append(tl.items, item.item)

	return nil
}

func (f *finder) List(limit storage.LimitFilter) ([]*model.Task, int, error) {
	listScanner := &listScanner{}
	totalCount, err := f.dbw.LoadList(f.builder, listScanner, limit)

	return listScanner.items, totalCount, err
}
