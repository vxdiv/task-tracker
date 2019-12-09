package sqldb

import (
	"database/sql"
	"time"

	. "github.com/go-xorm/builder"

	"github.com/vxdiv/task-tracker/app"
)

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{NewDBWrapper(db, "tasks")}
}

type TaskRepo struct {
	dbw *Wrapper
}

var _ app.TaskRepo = &TaskRepo{}

func (r *TaskRepo) Create(task *app.Task) error {
	var err error
	task.CreatedAt = time.Now().UTC()
	task.UpdatedAt = task.CreatedAt
	task.ID, err = r.dbw.Insert(Eq{
		"name":        task.Name,
		"description": task.Description,
		"type":        task.Type,
		"owner_id":    task.OwnerID,
		"assigned_id": task.AssignedID,
		"status":      task.Status,
		"due_date":    task.DueDate,
		"resolution":  task.Resolution,
		"priority":    task.Priority,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
	})

	if IsDuplicateEntryError(err) {
		return app.ErrTaskAlreadyExists
	}

	return err
}

func (r *TaskRepo) Update(task *app.Task) error {
	task.UpdatedAt = time.Now().UTC()
	return r.dbw.Update(Eq{
		"name":        task.Name,
		"description": task.Description,
		"type":        task.Type,
		"owner_id":    task.OwnerID,
		"assigned_id": task.AssignedID,
		"status":      task.Status,
		"due_date":    task.DueDate,
		"resolution":  task.Resolution,
		"priority":    task.Priority,
		"updated_at":  task.UpdatedAt,
	}, Eq{"id": task.ID})
}

func (r *TaskRepo) selectBuilder() *Builder {
	return r.dbw.SelectBuilder(
		"id",
		"name",
		"description",
		"type",
		"owner_id",
		"assigned_id",
		"status",
		"due_date",
		"resolution",
		"priority",
		"created_at",
		"updated_at",
	)
}

func (r *TaskRepo) scanTo(m *app.Task) (dst []interface{}) {
	return []interface{}{
		&m.ID,
		&m.Name,
		&m.Description,
		&m.Type,
		&m.OwnerID,
		&m.AssignedID,
		&m.Status,
		&m.DueDate,
		&m.Resolution,
		&m.Priority,
		&m.CreatedAt,
		&m.UpdatedAt,
	}
}

func (r *TaskRepo) GetByID(id int64) (*app.Task, error) {
	b := r.selectBuilder().Where(Eq{"id": id})

	return r.one(b)
}

func (r *TaskRepo) filter(b *Builder, filter app.TaskFilter) {
	if filter.CreatedAt != nil {
		filterTime(b, filter.CreatedAt, "created_at")
	}
}

func (r *TaskRepo) Count(f app.TaskFilter) (int, error) {
	b := r.selectBuilder()
	r.filter(b, f)

	return r.dbw.Count(b)
}

func (r *TaskRepo) One(f app.TaskFilter) (*app.Task, error) {
	b := r.selectBuilder()
	r.filter(b, f)

	return r.one(b)
}

func (r *TaskRepo) one(b *Builder) (*app.Task, error) {
	user := &app.Task{}
	err := r.dbw.Row(b, func(s RowScanner) error {
		return s.Scan(r.scanTo(user)...)
	})
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, app.ErrTaskNotFound
	default:
		return nil, err
	}

	return user, nil
}

func (r *TaskRepo) List(f app.TaskFilter) ([]*app.Task, error) {
	b := r.selectBuilder()
	r.filter(b, f)
	pagination(b, f.Pagination)

	var list []*app.Task
	err := r.dbw.Rows(b, func(s RowScanner) error {
		m := &app.Task{}
		list = append(list, m)
		return s.Scan(r.scanTo(m)...)
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}
