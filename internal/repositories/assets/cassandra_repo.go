package assets_db

import (
	assets_dm "assets/internal/core/domain/assets"
	"assets/internal/core/ports"
	"assets/pkg/logging"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"time"
)

const cassandraMaxLimit = 10_000

type CassandraRepo struct {
	logger        logging.Logger
	session       *gocql.Session
	chartsRepo    ports.ChartsRepository
	insightsRepo  ports.InsightsRepository
	audiencesRepo ports.AudiencesRepository
}

func NewCassandraRepo(logger logging.Logger, session *gocql.Session, chartsRepo ports.ChartsRepository, insightsRepo ports.InsightsRepository, audiencesRepo ports.AudiencesRepository) (repo *CassandraRepo) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := session.Query(CreateTableQuery()).WithContext(ctx).Exec(); err != nil {
		panic(errors.Wrap(err, "failed to inspect/create assets table"))
	}

	return &CassandraRepo{logger: logger, session: session, chartsRepo: chartsRepo, insightsRepo: insightsRepo, audiencesRepo: audiencesRepo}
}

func (cr *CassandraRepo) Select(ctx context.Context, params ports.SelectAssetsRepoParams) (results []assets_dm.AssetEntity, next string, err error) {

	cr.logger.Info("assets_db.Select() performed",
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

	var obj assets_dm.AssetEntity

	scanner := iter.Scanner()
	for scanner.Next() {
		if err = scanner.Scan(&obj.Id, &obj.ContentId, &obj.CreateTime, &obj.Description, &obj.Name, &obj.Type, &obj.UpdateTime); err != nil {
			fmt.Println("xdxd", err)
			return nil, next, err
		} else {
			results = append(results, obj)
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("xdxd2", err)

		return nil, next, err
	}

	if results, err = cr.populate(ctx, results...); err != nil {
		fmt.Println("xdxd3", err)

		return nil, next, err
	}

	return results, next, nil
}

func (cr *CassandraRepo) Insert(ctx context.Context, assets ...assets_dm.AssetEntity) (results []assets_dm.AssetEntity, err error) {

	cr.logger.Info("assets_db.Insert() performed",
		"params", assets,
		"results", results,
	)

	if err = cr.execute(ctx, assets, AppendInsertQuery); err != nil {
		return nil, err
	}

	if results, err = cr.populate(ctx, assets...); err != nil {
		return nil, err
	}

	return results, nil
}

func (cr *CassandraRepo) Update(ctx context.Context, models ...assets_dm.AssetEntity) (results []assets_dm.AssetEntity, err error) {

	cr.logger.Info("assets_db.Update() performed",
		"params", models,
		"results", results,
	)

	if len(models) == 0 {
		return results, nil
	}

	if err = cr.execute(ctx, models, AppendUpdateQuery); err != nil {
		return nil, err
	}

	if results, err = cr.populate(ctx, models...); err != nil {
		return nil, err
	}

	return results, nil
}

func (cr *CassandraRepo) Delete(ctx context.Context, models ...assets_dm.AssetEntity) (results []assets_dm.AssetEntity, err error) {

	cr.logger.Info("assets_db.Delete() performed",
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

func (cr *CassandraRepo) execute(ctx context.Context, assets []assets_dm.AssetEntity, action func(batch *gocql.Batch, asset assets_dm.AssetEntity)) (err error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	batch := cr.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	for idx := range assets {
		assets[idx].UpdateTime = time.Now()
		action(batch, assets[idx])
	}

	if err = cr.session.ExecuteBatch(batch); err != nil {
		return err
	}

	return nil
}
