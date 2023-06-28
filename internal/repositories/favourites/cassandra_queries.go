package favourites_db

import (
	favourites_dm "assets/internal/core/domain/favourites"
	"fmt"
	"github.com/gocql/gocql"
	"strings"
)

/*
 * Select
 */

var tableName = "favourites"

func SelectRecords(session *gocql.Session) (query *gocql.Query) {
	return session.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
}

func SelectRecordsByIds(session *gocql.Session, ids []string) (query *gocql.Query) {
	idList := "'" + strings.Join(ids, "', '") + "'"
	return session.Query(fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", tableName, idList))
}

func SelectRecordsByUserIds(session *gocql.Session, userIds []string) (query *gocql.Query) {
	idList := "'" + strings.Join(userIds, "', '") + "'"
	return session.Query(fmt.Sprintf("SELECT * FROM %s WHERE user_id IN (%s)", tableName, idList))
}

func SelectRecordsByAssetIds(session *gocql.Session, assetIds []string) (query *gocql.Query) {
	idList := "'" + strings.Join(assetIds, "', '") + "'"
	return session.Query(fmt.Sprintf("SELECT * FROM %s WHERE asset_id IN (%s)", tableName, idList))
}

/*
 * Table
 */

func CreateTableQuery() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id text PRIMARY KEY, user_id text, asset_id text, create_time timestamp, update_time timestamp)", tableName)
}

func CreateUserIdIndexQuery() string {
	return fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_user_id ON %s (user_id);", tableName)
}

func CreateAssetIdIndexQuery() string {
	return fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_asset_id ON %s (asset_id);", tableName)
}

func DropTableQuery() string {
	return fmt.Sprintf("DROP TABLE %s", tableName)
}

/*
 * Insert
 */

func AppendInsertQuery(batch *gocql.Batch, obj favourites_dm.FavouriteEntity) {
	batch.Query(fmt.Sprintf("INSERT INTO %s (id, user_id, asset_id, create_time, update_time) VALUES (?, ?, ?, ?, ?)", tableName),
		obj.Id, obj.UserId, obj.AssetId, obj.CreateTime, obj.UpdateTime)
}

/*
 * Delete
 */

func AppendDeleteQuery(batch *gocql.Batch, obj favourites_dm.FavouriteEntity) {
	batch.Query(fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName), obj.Id)
}
