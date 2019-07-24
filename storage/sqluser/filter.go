package sqluser

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

func (f *Filter) Name(name string) *Filter {
	f.builder.Where(Eq{"name": name})

	return f
}

func (f *Filter) Email(email string) *Filter {
	f.builder.Where(Eq{"email": email})

	return f
}

func (f *Filter) CreatedAt(filter sqldbw.TimeFilter) *Filter {
	if filter != nil {
		sqldbw.FilterTime(f.builder, filter, "created_at")
	}

	return f
}

func (f *Filter) Status(status model.UserStatus) *Filter {
	if len(status) > 0 {
		f.builder.Where(Eq{"status": status})
	}

	return f
}

func (f *Filter) Limit(limit sqldbw.Limit) *Filter {
	if limit != nil {
		f.builder.Limit(limit.Count(), limit.Offset())
	}

	return f
}

func (f *Filter) One() (*model.User, error) {
	return f.repo.One(f)
}

func (f *Filter) List() ([]*model.User, error) {
	return f.repo.List(f)
}

func (f *Filter) Count() (int, error) {
	return f.repo.Count(f)
}
