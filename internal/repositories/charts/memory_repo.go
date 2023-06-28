package charts_db

import (
	assets_dm "assets/internal/core/domain/assets"
	"assets/internal/core/ports"
	errs "assets/pkg/errors"
	"context"
)

/// test purposes database

type InMemoryDb struct {
	data map[string]assets_dm.ChartEntity
}

func NewMemoryRepo() *InMemoryDb {
	return &InMemoryDb{
		data: make(map[string]assets_dm.ChartEntity),
	}
}

func (i *InMemoryDb) Select(_ context.Context, params ports.SelectChartsRepoParams) (results []assets_dm.ChartEntity, cursor string, err error) {

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

	return results, cursor, err
}

func (i *InMemoryDb) Insert(_ context.Context, models ...assets_dm.ChartEntity) (results []assets_dm.ChartEntity, err error) {
	for _, model := range models {
		if _, ok := i.data[model.Id]; ok {
			return []assets_dm.ChartEntity{}, errs.AlreadyExistsError
		}
	}

	for _, model := range models {
		i.data[model.Id] = model
	}

	return models, err
}

func (i *InMemoryDb) Update(_ context.Context, models ...assets_dm.ChartEntity) (results []assets_dm.ChartEntity, err error) {

	for _, model := range models {
		if _, ok := i.data[model.Id]; !ok {
			return []assets_dm.ChartEntity{}, errs.CannotBeFoundError
		} else {
			i.data[model.Id] = model
		}
	}

	return models, err
}

func (i *InMemoryDb) Delete(_ context.Context, models ...assets_dm.ChartEntity) (results []assets_dm.ChartEntity, err error) {

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
