package users_db

import (
	users_dm "assets/internal/core/domain/users"
	"assets/internal/core/ports"
	errs "assets/pkg/errors"
	"context"
)

/// test purposes database

type InMemoryDb struct {
	data map[string]users_dm.UserEntity
}

func NewMemoryRepo() *InMemoryDb {
	return &InMemoryDb{
		data: make(map[string]users_dm.UserEntity),
	}
}

func (i *InMemoryDb) Select(_ context.Context, params ports.SelectUsersRepoParams) (results []users_dm.UserEntity, cursor string, err error) {

	if len(params.Ids) != 0 {
		for _, id := range params.Ids {
			if value, ok := i.data[id]; ok {
				results = append(results, value)
			}
		}
	}

	if len(results) > 0 {
		return results, cursor, err
	}

	if len(params.Emails) != 0 {
		for _, model := range i.data {
			for _, email := range params.Emails {
				if email == model.Email {
					results = append(results, model)
				}
			}
		}
	}

	if len(results) > 0 {
		return results, cursor, err
	}

	return results, cursor, err
}

func (i *InMemoryDb) Insert(_ context.Context, models ...users_dm.UserEntity) (results []users_dm.UserEntity, err error) {
	for _, model := range models {
		if _, ok := i.data[model.Id]; ok {
			return []users_dm.UserEntity{}, errs.AlreadyExistsError
		}
	}

	for _, model := range models {
		i.data[model.Id] = model
	}

	return models, err
}
