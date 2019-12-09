package sqldb

import (
	"github.com/go-xorm/builder"

	"github.com/vxdiv/task-tracker/app"
)

func filterTime(b *builder.Builder, filter app.TimeFilter, fieldName string) {
	if filter.From().After(filter.To()) {
		b.Where(builder.Gte{fieldName: filter.From()})
	} else {
		b.Where(builder.Gte{fieldName: filter.From()}.And(builder.Lte{fieldName: filter.To()}))
	}
}

func pagination(b *builder.Builder, filter app.Pager) {
	b.Limit(filter.Count(), filter.Offset())
}
