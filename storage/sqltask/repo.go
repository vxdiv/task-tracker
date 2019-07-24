package sqltask

import (
	"database/sql"
	"errors"
	"time"

	. "github.com/go-xorm/builder"

	"github.com/vxdiv/task-tracker/model"
	"github.com/vxdiv/task-tracker/storage/sqldbw"
)

var (
	ErrTaskNotFound      = errors.New("task is not found")
	ErrTaskAlreadyExists = errors.New("task is already exist")
)

type Repo interface {
	Create(task *model.Task) error
	Update(task *model.Task) error
	Filter() *Filter
	Count(f *Filter) (int, error)
	One(f *Filter) (*model.Task, error)
	List(f *Filter) ([]*model.Task, error)
}

func New(db *sql.DB) *SqlRepo {
	return &SqlRepo{sqldbw.New(db, "tasks")}
}

type SqlRepo struct {
	dbw *sqldbw.Wrapper
}

var _ Repo = &SqlRepo{}

func (r *SqlRepo) Create(task *model.Task) error {
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

	if sqldbw.IsDuplicateEntryError(err) {
		return ErrTaskAlreadyExists
	}

	return err
}

func (r *SqlRepo) Update(task *model.Task) error {
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

func (r *SqlRepo) Filter() *Filter {
	return &Filter{
		repo: r,
		builder: r.dbw.SelectBuilder(
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
		)}
}

func scanTo(m *model.Task) (dst []interface{}) {
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

func (r *SqlRepo) Count(f *Filter) (int, error) {
	return r.dbw.Count(f.builder)
}

func (r *SqlRepo) One(f *Filter) (*model.Task, error) {
	m := &model.Task{}
	err := r.dbw.Row(f.builder, func(s sqldbw.RowScanner) error {
		return s.Scan(scanTo(m)...)
	})
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrTaskNotFound
	default:
		return nil, err
	}

	return m, nil
}

func (r *SqlRepo) List(f *Filter) ([]*model.Task, error) {
	var list []*model.Task
	err := r.dbw.Rows(f.builder, func(s sqldbw.RowScanner) error {
		m := &model.Task{}
		list = append(list, m)
		return s.Scan(scanTo(m)...)
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}
