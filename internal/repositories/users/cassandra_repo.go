package users_db

import (
	users_dm "assets/internal/core/domain/users"
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
		panic(errors.Wrap(err, "failed to inspect/create users table"))
	}

	if err := session.Query(CreateEmailIndexQuery()).WithContext(ctx).Exec(); err != nil {
		panic(errors.Wrap(err, "failed to inspect/create email index on messages table"))
	}

	return &CassandraRepo{logger: logger, session: session}
}

func (cr *CassandraRepo) Select(ctx context.Context, params ports.SelectUsersRepoParams) (results []users_dm.UserEntity, next string, err error) {

	cr.logger.Info("users_db.Select() performed",
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
	} else if len(params.Emails) != 0 {
		query = SelectRecordsByEmails(cr.session, params.Emails)
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

	var obj users_dm.UserEntity

	scanner := iter.Scanner()
	for scanner.Next() {
		if err = scanner.Scan(&obj.Id, &obj.CreateTime, &obj.Email, &obj.Password, &obj.UpdateTime); err != nil {
			return nil, next, err
		} else {
			results = append(results, obj)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, next, err
	}

	return results, next, nil
}

func (cr *CassandraRepo) Insert(ctx context.Context, assets ...users_dm.UserEntity) (results []users_dm.UserEntity, err error) {

	cr.logger.Info("users_db.Insert() performed",
		"params", assets,
		"results", results,
	)

	if err = cr.execute(ctx, assets, AppendInsertQuery); err != nil {
		return nil, err
	}

	return assets, nil
}

func (cr *CassandraRepo) execute(ctx context.Context, users []users_dm.UserEntity, action func(batch *gocql.Batch, user users_dm.UserEntity)) (err error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	batch := cr.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	for idx := range users {
		users[idx].UpdateTime = time.Now()
		action(batch, users[idx])
	}

	if err = cr.session.ExecuteBatch(batch); err != nil {
		return err
	}

	return nil
}
