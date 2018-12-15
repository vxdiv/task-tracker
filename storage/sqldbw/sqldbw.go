package sqldbw

import (
	"database/sql"

	. "github.com/go-xorm/builder"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "db")

func SetLogger(l *logrus.Entry) {
	log = l
}

const dbDialect = MYSQL

type Scanner interface {
	Scan(dst ...interface{}) error
}

type ItemScanner interface {
	ScanItem(s Scanner) error
}

type LimitFilter interface {
	Offset() int
	Count() int
}

type Wrapper struct {
	db    *sql.DB
	table string
}

func New(db *sql.DB, table string) *Wrapper {
	return &Wrapper{db: db, table: table}
}

func (rw *Wrapper) Insert(eq Eq) (int64, error) {
	sqlQuery, args, err := Dialect(dbDialect).Insert(eq).Into(rw.table).ToSQL()
	if err != nil {
		return 0, err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	res, err := rw.db.Exec(sqlQuery, args...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (rw *Wrapper) Update(eq Eq, cond Cond) error {
	sqlQuery, args, err := Update(eq).From(rw.table).Where(cond).ToSQL()
	if err != nil {
		return err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	_, err = rw.db.Exec(sqlQuery, args...)

	return err
}

func (rw *Wrapper) SelectBuilder(cols ...string) *Builder {
	return Dialect(dbDialect).Select(cols...).From(rw.table)
}

func (rw *Wrapper) LoadOne(b *Builder, s ItemScanner) error {
	sqlQuery, args, err := b.ToSQL()
	if err != nil {
		return err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	row := rw.db.QueryRow(sqlQuery, args...)

	return s.ScanItem(row)
}

func (rw *Wrapper) LoadList(b *Builder, s ItemScanner, limitFilter LimitFilter) (totalRowCount int, err error) {
	var (
		countQuery string
		sqlQuery   string
		args       []interface{}
	)

	bCount := *b
	countQuery, args, err = bCount.Select("count(*)").ToSQL()
	if err != nil {
		return 0, err
	}

	log.Debugf("sql: %s; with args: %+v", countQuery, args)

	err = rw.db.QueryRow(countQuery, args...).Scan(&totalRowCount)
	if err != nil {
		return 0, err
	}

	if limitFilter != nil {
		b.Limit(limitFilter.Count(), limitFilter.Offset())
	}

	sqlQuery, args, err = b.ToSQL()
	if err != nil {
		return 0, err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	rows, err := rw.db.Query(sqlQuery, args...)
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := s.ScanItem(rows); err != nil {
			return 0, err
		}
	}

	return totalRowCount, rows.Err()
}
