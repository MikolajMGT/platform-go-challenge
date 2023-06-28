package audiences_db

import (
	assets_dm "assets/internal/core/domain/assets"
	"fmt"
	"github.com/gocql/gocql"
	"strings"
)

/*
 * Select
 */

var tableName = "audiences"

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
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id text PRIMARY KEY, gender text, birth_country text, age_group text, social_media_hours int, purchases_last_month int, create_time timestamp, update_time timestamp)", tableName)
}

func DropTableQuery() string {
	return fmt.Sprintf("DROP TABLE %s", tableName)
}

/*
 * Insert
 */

func AppendInsertQuery(batch *gocql.Batch, obj assets_dm.AudienceEntity) {
	batch.Query(fmt.Sprintf("INSERT INTO %s (id, gender, birth_country, age_group, social_media_hours, purchases_last_month, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", tableName),
		obj.Id, obj.Gender, obj.BirthCountry, obj.AgeGroup, obj.SocialMediaHours, obj.PurchasesLastMonth, obj.CreateTime, obj.UpdateTime)
}

/*
 * Delete
 */

func AppendDeleteQuery(batch *gocql.Batch, obj assets_dm.AudienceEntity) {
	batch.Query(fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName), obj.Id)
}
