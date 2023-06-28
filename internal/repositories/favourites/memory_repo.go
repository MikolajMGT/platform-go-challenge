package favourites_db

import (
	favourites "assets/internal/core/domain/favourites"
	"assets/internal/core/ports"
	errs "assets/pkg/errors"
	"context"
)

/// test purposes database

type InMemoryDb struct {
	data map[string]favourites.FavouriteEntity
}

func NewMemoryRepo() *InMemoryDb {
	return &InMemoryDb{
		data: make(map[string]favourites.FavouriteEntity),
	}
}

func (i *InMemoryDb) Select(_ context.Context, params ports.SelectFavouritesRepoParams) (results []favourites.FavouriteEntity, cursor string, err error) {

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

	if len(params.UserIds) != 0 {
		for _, model := range i.data {
			for _, userId := range params.UserIds {
				if userId == model.UserId {
					results = append(results, model)
				}
			}
		}
	}

	if len(results) > 0 {
		return results, cursor, err
	}

	if len(params.AssetIds) != 0 {
		for _, model := range i.data {
			for _, assetId := range params.AssetIds {
				if assetId == model.AssetId {
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

func (i *InMemoryDb) Insert(_ context.Context, models ...favourites.FavouriteEntity) (results []favourites.FavouriteEntity, err error) {
	for _, model := range models {
		if _, ok := i.data[model.Id]; ok {
			return []favourites.FavouriteEntity{}, errs.AlreadyExistsError
		}
	}

	for _, model := range models {
		i.data[model.Id] = model
	}

	return models, err
}

func (i *InMemoryDb) Update(_ context.Context, models ...favourites.FavouriteEntity) (results []favourites.FavouriteEntity, err error) {

	for _, model := range models {
		if _, ok := i.data[model.Id]; !ok {
			return []favourites.FavouriteEntity{}, errs.CannotBeFoundError
		} else {
			i.data[model.Id] = model
		}
	}

	return models, err
}

func (i *InMemoryDb) Delete(_ context.Context, models ...favourites.FavouriteEntity) (results []favourites.FavouriteEntity, err error) {

	for _, model := range models {
		if _, ok := i.data[model.Id]; !ok {
			return nil, errs.CannotBeFoundError
		} else {
			results = append(results, model)
			delete(i.data, model.Id)
		}
	}

	return results, nil
}
