package sessions_db

import (
	"fmt"
	"github.com/gocql/gocql"
)

/*
 * Select
 */

var tableName = "sessions"

func SelectRecordsById(session *gocql.Session, id string) (query *gocql.Query) {
	return session.Query(fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tableName, id))
}

/*
 * Table
 */

func CreateTableQuery() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id text PRIMARY KEY, \"values\" text, create_time timestamp, update_time timestamp)", tableName)
}

func DropTableQuery() string {
	return fmt.Sprintf("DROP TABLE %s", tableName)
}

/*
 * Insert
 */

func InsertQuery(cassandra *gocql.Session, id string, values string) (query *gocql.Query) {
	return cassandra.Query(fmt.Sprintf("INSERT INTO %s (id, \"values\") VALUES (?, ?)", tableName),
		id, values)
}
