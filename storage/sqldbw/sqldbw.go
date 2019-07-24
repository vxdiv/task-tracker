package sqldbw

import (
	"database/sql"

	. "github.com/go-xorm/builder"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "db")

const dbDialect = MYSQL

type RowScanner interface {
	Scan(dst ...interface{}) error
}

type ScanFunc func(RowScanner) error

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
	sqlQuery, args, err := Dialect(dbDialect).Update(eq).From(rw.table).Where(cond).ToSQL()
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

func (rw *Wrapper) Row(b *Builder, scanner ScanFunc) error {
	sqlQuery, args, err := b.ToSQL()
	if err != nil {
		return err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	row := rw.db.QueryRow(sqlQuery, args...)

	return scanner(row)
}

func (rw *Wrapper) Count(b *Builder) (int, error) {
	countB := *b
	countQuery, args, err := countB.Select("count(*)").ToSQL()
	if err != nil {
		return 0, err
	}

	log.Debugf("sql: %s; with args: %+v", countQuery, args)

	var totalCount int
	err = rw.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (rw *Wrapper) Rows(b *Builder, scanner ScanFunc) error {

	sqlQuery, args, err := b.ToSQL()
	if err != nil {
		return err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	rows, err := rw.db.Query(sqlQuery, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := scanner(rows); err != nil {
			return err
		}
	}

	return rows.Err()
}
