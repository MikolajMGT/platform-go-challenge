package charts_db

import (
	assets_dm "assets/internal/core/domain/assets"
	"assets/internal/core/ports"
	"assets/pkg/logging"
	"context"
	"encoding/base64"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"time"
)

const cassandraMaxLimit = 10_000

type CassandraRepo struct {
	logger  logging.Logger
	session *gocql.Session
}

func NewCassandraRepo(logger logging.Logger, session *gocql.Session) (repo *CassandraRepo) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := session.Query(CreateTableQuery()).WithContext(ctx).Exec(); err != nil {
		panic(errors.Wrap(err, "failed to inspect/create charts table"))
	}

	return &CassandraRepo{logger: logger, session: session}
}

func (cr *CassandraRepo) Select(ctx context.Context, params ports.SelectChartsRepoParams) (results []assets_dm.ChartEntity, next string, err error) {

	cr.logger.Info("charts_db.Select() performed",
		"params", params,
		"results", results,
	)

	var (
		limit  = cassandraMaxLimit
		cursor = make([]byte, 0)
	)

	if params.Limit != 0 {
		limit = params.Limit
	}

	if params.Cursor != "" {
		if cursor, err = base64.URLEncoding.DecodeString(params.Cursor); err != nil {
			return nil, next, err
		}
	}

	var query *gocql.Query
	if len(params.Ids) != 0 {
		query = SelectRecordsByIds(cr.session, params.Ids)
	} else {
		query = SelectRecords(cr.session)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	iter := query.WithContext(ctx).PageSize(limit).PageState(cursor).Iter()
	defer func() {
		if err = iter.Close(); err != nil {
			cr.logger.Info("failed to close iterator", "err", err)
		}
	}()

	if len(iter.PageState()) > 0 {
		next = base64.URLEncoding.EncodeToString(iter.PageState())
	}

	var obj assets_dm.ChartEntity

	scanner := iter.Scanner()
	for scanner.Next() {
		var data string
		if err = scanner.Scan(&obj.Id, &obj.ChartTitle, &obj.CreateTime, &data, &obj.UpdateTime, &obj.XAxisTitle, &obj.YAxisTitle); err != nil {
			return nil, next, err
		} else {
			obj.Data = data
			results = append(results, obj)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, next, err
	}

	return results, next, nil
}

func (cr *CassandraRepo) Insert(ctx context.Context, models ...assets_dm.ChartEntity) (results []assets_dm.ChartEntity, err error) {

	cr.logger.Info("charts_db.Insert() performed",
		"params", models,
		"results", results,
	)

	if len(models) == 0 {
		return results, nil
	}

	if err = cr.execute(ctx, models, AppendInsertQuery); err != nil {
		return nil, err
	}

	return models, nil
}

func (cr *CassandraRepo) Delete(ctx context.Context, models ...assets_dm.ChartEntity) (results []assets_dm.ChartEntity, err error) {

	cr.logger.Info("charts_db.Delete() performed",
		"params", models,
		"results", results,
	)

	if len(models) == 0 {
		return results, nil
	}

	if err = cr.execute(ctx, models, AppendDeleteQuery); err != nil {
		return nil, err
	}

	return models, nil
}

func (cr *CassandraRepo) execute(ctx context.Context, models []assets_dm.ChartEntity, action func(batch *gocql.Batch, asset assets_dm.ChartEntity)) (err error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	batch := cr.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	for idx := range models {
		models[idx].UpdateTime = time.Now()
		action(batch, models[idx])
	}

	if err = cr.session.ExecuteBatch(batch); err != nil {
		return err
	}

	return nil
}
