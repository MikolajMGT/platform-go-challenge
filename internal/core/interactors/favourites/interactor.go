package favourites_itc

import (
	assets_dm "assets/internal/core/domain/assets"
	favourites_dm "assets/internal/core/domain/favourites"
	users_dm "assets/internal/core/domain/users"
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
	favouritesRepo ports.FavouritesRepository
	usersRepo      ports.UsersRepository
	assetsRepo     ports.AssetsRepository
}

func NewInteractor(logger logging.Logger, validator validation.Validator, favouritesRepo ports.FavouritesRepository, usersRepo ports.UsersRepository, assetsRepo ports.AssetsRepository) *Interactor {
	return &Interactor{
		logger:         logger,
		validator:      validator,
		favouritesRepo: favouritesRepo,
		usersRepo:      usersRepo,
		assetsRepo:     assetsRepo,
	}
}

func (i *Interactor) Select(ctx context.Context, params ports.SelectFavouritesItcParams) (results []favourites_dm.FavouriteEntity, cursor string, err error) {

	i.logger.Info("favourites_itc.Select() performed",
		"params", params,
		"results", results,
	)

	if err = i.validator.Validate(params); err != nil {
		return nil, cursor, errors.Join(errs.ValidationError, err)
	}

	if results, cursor, err = i.favouritesRepo.Select(ctx, convertSelectParams(params)); err != nil {
		return nil, cursor, errors.Join(errs.ProcessingError, err)
	}

	return results, cursor, err
}

func (i *Interactor) Insert(ctx context.Context, params ...ports.InsertFavouriteItcParams) (results []favourites_dm.FavouriteEntity, err error) {

	i.logger.Info("favourites_itc.Insert() performed",
		"params", params,
		"results", results,
	)

	if err = i.validator.Validate(params); err != nil {
		return nil, errors.Join(errs.ValidationError, err)
	}

	userIds := slices.Map(params, func(param ports.InsertFavouriteItcParams) string {
		return param.UserId
	})

	var users []users_dm.UserEntity
	if users, _, err = i.usersRepo.Select(ctx, ports.SelectUsersRepoParams{Ids: userIds}); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	if len(userIds) != len(users) {
		return nil, errors.Join(errs.CannotBeFoundError, errors.New("user cannot be found"))
	}

	assetIds := slices.Map(params, func(param ports.InsertFavouriteItcParams) string {
		return param.AssetId
	})

	var assets []assets_dm.AssetEntity
	if assets, _, err = i.assetsRepo.Select(ctx, ports.SelectAssetsRepoParams{Ids: assetIds}); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	if len(assetIds) != len(assets) {
		return nil, errors.Join(errs.CannotBeFoundError, errors.New("asset cannot be found"))
	}

	var current []favourites_dm.FavouriteEntity
	if current, _, err = i.favouritesRepo.Select(ctx, ports.SelectFavouritesRepoParams{UserIds: userIds}); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	i.logger.Info("0000000--->", "current", current)

	currentAssetIds := slices.Map(current, func(obj favourites_dm.FavouriteEntity) string {
		return obj.AssetId
	})

	if slices.HasCommon(assetIds, currentAssetIds) {
		return nil, errors.Join(errs.AlreadyExistsError, errors.New("provided asset is already on favourites list"))
	}

	if results, err = i.favouritesRepo.Insert(ctx, prepareCreatableModels(params)...); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	return results, err
}

func (i *Interactor) Delete(ctx context.Context, params ...ports.DeleteFavouriteItcParams) (results []favourites_dm.FavouriteEntity, err error) {

	i.logger.Info("favourites_itc.Delete() performed",
		"params", params,
		"results", results,
	)

	if err = i.validator.Validate(params); err != nil {
		return nil, errors.Join(errs.ValidationError, err)
	}

	ids := slices.Map(params, func(param ports.DeleteFavouriteItcParams) string {
		return param.Id
	})

	var models []favourites_dm.FavouriteEntity
	if models, _, err = i.Select(ctx, ports.SelectFavouritesItcParams{
		Ids: ids,
	}); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	if results, err = i.favouritesRepo.Delete(ctx, models...); err != nil {
		return nil, errors.Join(errs.ProcessingError, err)
	}

	return results, err
}
