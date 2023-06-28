package users_db

import (
	users_dm "assets/internal/core/domain/users"
	"fmt"
	"github.com/gocql/gocql"
	"strings"
)

/*
 * Select
 */

var tableName = "users"

func SelectRecords(session *gocql.Session) (query *gocql.Query) {
	return session.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
}

func SelectRecordsByIds(session *gocql.Session, ids []string) (query *gocql.Query) {
	idList := "'" + strings.Join(ids, "', '") + "'"
	return session.Query(fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", tableName, idList))
}

func SelectRecordsByEmails(session *gocql.Session, emails []string) (query *gocql.Query) {
	emailsList := "'" + strings.Join(emails, "', '") + "'"
	return session.Query(fmt.Sprintf("SELECT * FROM %s WHERE email IN (%s)", tableName, emailsList))
}

/*
 * Table
 */

func CreateTableQuery() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id text PRIMARY KEY, email text, password text, create_time timestamp, update_time timestamp)", tableName)
}

func CreateEmailIndexQuery() string {
	return fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_email ON %s (email);", tableName)
}

func DropTableQuery() string {
	return fmt.Sprintf("DROP TABLE %s", tableName)
}

/*
 * Insert
 */

func AppendInsertQuery(batch *gocql.Batch, obj users_dm.UserEntity) {
	batch.Query(fmt.Sprintf("INSERT INTO %s (id, email, password, create_time, update_time) VALUES (?, ?, ?, ?, ?)", tableName),
		obj.Id, obj.Email, obj.Password, obj.CreateTime, obj.UpdateTime)
}
