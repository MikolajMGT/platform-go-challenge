package charts_db

import (
	assets_dm "assets/internal/core/domain/assets"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"strings"
)

/*
 * Select
 */

var tableName = "charts"

func SelectRecords(session *gocql.Session) (query *gocql.Query) {
	return session.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
}

func SelectRecordsByIds(session *gocql.Session, ids []string) (query *gocql.Query) {
	idList := "'" + strings.Join(ids, "', '") + "'"
	return session.Query(fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", tableName, idList))
}

/*
 * Table
 */

func CreateTableQuery() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id text PRIMARY KEY, chart_title text, x_axis_title text, y_axis_title text, \"data\" text, create_time timestamp, update_time timestamp)", tableName)
}

func DropTableQuery() string {
	return fmt.Sprintf("DROP TABLE %s", tableName)
}

/*
 * Insert
 */

func AppendInsertQuery(batch *gocql.Batch, obj assets_dm.ChartEntity) {
	data, _ := json.Marshal(obj.Data)
	batch.Query(fmt.Sprintf("INSERT INTO %s (id, chart_title, x_axis_title, y_axis_title, \"data\", create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?)", tableName),
		obj.Id, obj.ChartTitle, obj.XAxisTitle, obj.YAxisTitle, string(data), obj.CreateTime, obj.UpdateTime)
}

/*
 * Delete
 */

func AppendDeleteQuery(batch *gocql.Batch, obj assets_dm.ChartEntity) {
	batch.Query(fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName), obj.Id)
}
