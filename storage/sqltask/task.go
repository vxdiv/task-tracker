package sqltask

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

var _ storage.TaskRepo = &Repo{}

func New(db *sql.DB) *Repo {
	return &Repo{dbw: sqldbw.New(db, "tasks")}
}

func (r Repo) Create(task *model.Task) error {
	var err error
	task.CreatedAt = time.Now().UTC()
	task.UpdatedAt = task.CreatedAt
	task.ID, err = r.dbw.Insert(builder.Eq{
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
		return storage.ErrTaskAlreadyExists
	}

	return err
}

func (r Repo) Update(task *model.Task) error {
	task.UpdatedAt = time.Now().UTC()
	return r.dbw.Update(builder.Eq{
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
	}, builder.Eq{"id": task.ID})
}

func (r Repo) Find() storage.TaskFinder {
	return &finder{
		dbw: r.dbw,
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
