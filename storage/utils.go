package storage

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/builder"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "db")

const dbDialect = builder.MYSQL

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

type TimeFilter interface {
	From() time.Time
	To() time.Time
}

func isDuplicateEntryError(err error) bool {
	mysqlError, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	return mysqlError.Number == 1062
}

func insert(db *sql.DB, eq builder.Eq, table string) (int64, error) {
	sqlQuery, args, err := builder.Insert(eq).Into(table).ToSQL()
	if err != nil {
		return 0, err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	res, err := db.Exec(sqlQuery, args...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func update(db *sql.DB, eq, condEq builder.Eq, table string) error {
	sqlQuery, args, err := builder.Update(eq).From(table).Where(condEq).ToSQL()
	if err != nil {
		return err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	_, err = db.Exec(sqlQuery, args...)

	return err
}

func loadOne(db *sql.DB, b *builder.Builder, s ItemScanner) error {
	sqlQuery, args, err := b.ToSQL()
	if err != nil {
		return err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	row := db.QueryRow(sqlQuery, args...)

	return s.ScanItem(row)
}

func loadAll(db *sql.DB, b *builder.Builder, s ItemScanner, limitFilter LimitFilter) (int, error) {
	count := 0
	bCount := *b
	selectCount, args, err := bCount.Select("count(*)").ToSQL()
	if err != nil {
		return 0, err
	}

	log.Debugf("sql: %s; with args: %+v", selectCount, args)

	err = db.QueryRow(selectCount, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	if limitFilter != nil {
		b.Limit(limitFilter.Count(), limitFilter.Offset())
	}

	sqlQuery, args, err := b.ToSQL()
	if err != nil {
		return 0, err
	}

	log.Debugf("sql: %s; with args: %+v", sqlQuery, args)

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := s.ScanItem(rows); err != nil {
			return 0, err
		}
	}

	return count, rows.Err()
}

func filterTime(b *builder.Builder, filter TimeFilter, fieldName string) {
	if filter.From().After(filter.To()) {
		b.Where(builder.Gt{fieldName: filter.From()})
	} else {
		b.Where(
			builder.Gt{fieldName: filter.From()}.And(builder.Lt{fieldName: filter.To()}))
	}
}
