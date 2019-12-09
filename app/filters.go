package app

import (
	"time"
)

type TimeFilter interface {
	From() time.Time
	To() time.Time
}

type Pager interface {
	Offset() int
	Count() int
}
