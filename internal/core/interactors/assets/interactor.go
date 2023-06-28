package assets_itc

import (
	assets_dm "assets/internal/core/domain/assets"
	favourites_dm "assets/internal/core/domain/favourites"
	"assets/internal/core/ports"
	errs "assets/pkg/errors"
	"assets/pkg/logging"
	"assets/pkg/slices"
	"assets/pkg/validation"
	"context"
	"errors"
)

type Interactor struct {
	logger         logging.Logger
	validator      validation.Validator
	assetsRepo     ports.AssetsRepository
	chartsRepo     ports.ChartsRepository
	insightsRepo   ports.InsightsRepository
	audiencesRepo  ports.AudiencesRepository
	favouritesRepo ports.FavouritesRepository
}

func NewInteractor(logger logging.Logger, validator validation.Validator, assetsRepo ports.AssetsRepository, chartsRepo ports.ChartsRepository, insightsRepo ports.InsightsRepository, audiencesRepo ports.AudiencesRepository, favouritesRepo ports.FavouritesRepository) *Interactor {
	return &Interactor{
		logger:         logger,
		validator:      validator,
		assetsRepo:     assetsRepo,
		chartsRepo:     chartsRepo,
		insightsRepo:   insightsRepo,
		audiencesRepo:  audiencesRepo,
		favouritesRepo: favouritesRepo,
	}
}

func (i *Interactor) Select(ctx context.Context, params ports.SelectAssetsItcParams) (results []assets_dm.AssetEntity, cursor string, err error) {

	i.logger.Info("assets_itc.Select() performed",
		"params", params,
		"results", results,
	)

	if err = i.validator.Validate(params); err != nil {
		return nil, cursor, errors.Join(errs.ValidationError, err)
	}

	if results, cursor, err = i.assetsRepo.Select(ctx, convertSelectParams(params)); err != nil {
		return nil, cursor, errors.Join(errs.ProcessingError, err)
	}

	return results, cursor, err
}

func (i *Interactor) Insert(ctx context.Context, params ...ports.InsertAssetItcParams) (results []assets_dm.AssetEntity, err error) {

	i.logger.Info("assets_itc.Insert() performed",
		"params", params,
		"results", results,
	)

	if err = i.validator.Validate(params); err != nil {
		return nil, errors.Join(errs.ValidationError, err)
	}

	for _, param := range params {
		if param.AssetData.Chart == nil && param.AssetData.Insight == nil && param.AssetData.Audience == nil {
			return nil, errors.Join(errs.ValidationError, errors.New("asset data cannot be empty"))
		}
	}

	charts, insights, audiences, mapper, err := i.createDependencies(ctx, params...)
	if err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	defer func() {
		if err == nil {
			return
		}

		if _, err = i.chartsRepo.Delete(ctx, charts...); err != nil {
			i.logger.Info("failed to recreate data")
		}
		if _, err = i.insightsRepo.Delete(ctx, insights...); err != nil {
			i.logger.Info("failed to recreate data")
		}
		if _, err = i.audiencesRepo.Delete(ctx, audiences...); err != nil {
			i.logger.Info("failed to recreate data")
		}
	}()

	if results, err = prepareCreatableModels(params, mapper); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	if results, err = i.assetsRepo.Insert(ctx, results...); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	return results, err
}

func (i *Interactor) Update(ctx context.Context, params ...ports.UpdateAssetItcParams) (results []assets_dm.AssetEntity, err error) {

	i.logger.Info("assets_itc.Update() performed",
		"params", params,
		"results", results,
	)

	if err = i.validator.Validate(params); err != nil {
		return nil, errors.Join(errs.ValidationError, err)
	}

	ids := slices.Map(params, func(param ports.UpdateAssetItcParams) string {
		return param.Id
	})

	var models []assets_dm.AssetEntity
	if models, _, err = i.assetsRepo.Select(ctx, ports.SelectAssetsRepoParams{
		Ids: ids,
	}); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	if models, err = slices.MatchOrder(params, models, func(e1 ports.UpdateAssetItcParams, e2 assets_dm.AssetEntity) bool {
		return e1.Id == e2.Id
	}); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	if results, err = i.assetsRepo.Update(ctx, prepareUpdatableModels(params, models)...); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	return results, err
}

func (i *Interactor) Delete(ctx context.Context, params ...ports.DeleteAssetItcParams) (results []assets_dm.AssetEntity, err error) {

	i.logger.Info("assets_itc.Delete() performed",
		"params", params,
		"results", results,
	)

	if err = i.validator.Validate(params); err != nil {
		return nil, errors.Join(errs.ValidationError, err)
	}

	ids := slices.Map(params, func(param ports.DeleteAssetItcParams) string {
		return param.Id
	})

	var favourites []favourites_dm.FavouriteEntity
	if favourites, _, err = i.favouritesRepo.Select(ctx, ports.SelectFavouritesRepoParams{AssetIds: ids}); err != nil {
		return nil, err
	}

	if _, err = i.favouritesRepo.Delete(ctx, favourites...); err != nil {
		return nil, err
	}

	defer func() {
		if err == nil {
			return
		}

		if _, err = i.favouritesRepo.Insert(ctx, favourites...); err != nil {
			i.logger.Info("failed to recreate data")
		}
	}()

	var models []assets_dm.AssetEntity
	if models, _, err = i.Select(ctx, ports.SelectAssetsItcParams{
		Ids: ids,
	}); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	charts, insights, audiences, err := i.deleteDependencies(ctx, models...)
	if err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	defer func() {
		if err == nil {
			return
		}

		if _, err = i.chartsRepo.Insert(ctx, charts...); err != nil {
			i.logger.Info("failed to recreate data")
		}
		if _, err = i.insightsRepo.Insert(ctx, insights...); err != nil {
			i.logger.Info("failed to recreate data")
		}
		if _, err = i.audiencesRepo.Insert(ctx, audiences...); err != nil {
			i.logger.Info("failed to recreate data")
		}
	}()

	if results, err = i.assetsRepo.Delete(ctx, models...); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	return results, err
}
