package sqltask

import (
	. "github.com/go-xorm/builder"

	"github.com/vxdiv/task-tracker/model"
	"github.com/vxdiv/task-tracker/storage/sqldbw"
)

type Filter struct {
	repo    Repo
	builder *Builder
}

func (f *Filter) ID(id int64) *Filter {
	f.builder.Where(Eq{"id": id})

	return f
}

func (f *Filter) Status(status string) *Filter {
	if len(status) > 0 {
		f.builder.Where(Eq{"status": status})
	}

	return f
}

func (f *Filter) Type(typ string) *Filter {
	if len(typ) > 0 {
		f.builder.Where(Eq{"type": typ})
	}

	return f
}

func (f *Filter) OwnerID(id int64) *Filter {
	if id > 0 {
		f.builder.Where(Eq{"owner_id": id})
	}

	return f
}

func (f *Filter) AssignedID(id int64) *Filter {
	if id > 0 {
		f.builder.Where(Eq{"assigned_id": id})
	}

	return f
}

func (f *Filter) Resolution(resolution string) *Filter {
	if len(resolution) > 0 {
		f.builder.Where(Eq{"resolution": resolution})
	}

	return f
}

func (f *Filter) Priority(priority string) *Filter {
	if len(priority) > 0 {
		f.builder.Where(Eq{"priority": priority})
	}

	return f
}

func (f *Filter) NameLike(name string) *Filter {
	if len(name) > 0 {
		f.builder.Where(Like{"name", "%" + name})
	}

	return f
}

func (f *Filter) CreatedAt(filter sqldbw.TimeFilter) *Filter {
	if filter != nil {
		sqldbw.FilterTime(f.builder, filter, "created_at")
	}

	return f
}

func (f *Filter) DueDate(filter sqldbw.TimeFilter) *Filter {
	if filter != nil {
		sqldbw.FilterTime(f.builder, filter, "due_date")
	}

	return f
}

func (f *Filter) Limit(limit sqldbw.Limit) *Filter {
	if limit != nil {
		f.builder.Limit(limit.Count(), limit.Offset())
	}

	return f
}

func (f *Filter) One() (*model.Task, error) {
	return f.repo.One(f)
}

func (f *Filter) List() ([]*model.Task, error) {
	return f.repo.List(f)
}

func (f *Filter) Count() (int, error) {
	return f.repo.Count(f)
}
