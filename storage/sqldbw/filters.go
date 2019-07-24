package sqldbw

import (
	"time"

	. "github.com/go-xorm/builder"
)

type Limit interface {
	Offset() int
	Count() int
}

type TimeFilter interface {
	From() time.Time
	To() time.Time
}

func FilterTime(b *Builder, filter TimeFilter, fieldName string) {
	if filter.From().After(filter.To()) {
		b.Where(Gte{fieldName: filter.From()})
	} else {
		b.Where(Gte{fieldName: filter.From()}.And(Lte{fieldName: filter.To()}))
	}
}
