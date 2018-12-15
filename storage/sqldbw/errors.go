package sqldbw

import "github.com/go-sql-driver/mysql"

func IsDuplicateEntryError(err error) bool {
	mysqlError, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	return mysqlError.Number == 1062
}
