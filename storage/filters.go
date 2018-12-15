package storage

import (
	"time"

	"github.com/go-xorm/builder"
	"github.com/vxdiv/task-tracker/storage/sqldbw"
)

type LimitFilter sqldbw.LimitFilter

type TimeFilter interface {
	From() time.Time
	To() time.Time
}

func FilterTime(b *builder.Builder, filter TimeFilter, fieldName string) {
	if filter.From().After(filter.To()) {
		b.Where(builder.Gt{fieldName: filter.From()})
	} else {
		b.Where(
			builder.Gt{fieldName: filter.From()}.And(builder.Lt{fieldName: filter.To()}))
	}
}
