package ports

import (
	assets_dm "assets/internal/core/domain/assets"
	favourites_dm "assets/internal/core/domain/favourites"
	users_dm "assets/internal/core/domain/users"
	"context"
)

/*
 * Users
 */

/// params

type LoginUserItcParams struct {
	Email    string `validate:"required,email,max=64" json:"email"`
	Password string `validate:"required,max=64" json:"password"`
}

type RegisterUserItcParams struct {
	Email    string `validate:"required,email,max=64" json:"email"`
	Password string `validate:"required,max=64" json:"password"`
}

/// interactor

type UsersInteractor interface {
	Login(ctx context.Context, params LoginUserItcParams) (users_dm.UserEntity, error)
	Register(ctx context.Context, params RegisterUserItcParams) (users_dm.UserEntity, error)
}

/*
 * Assets
 */

/// params

type SelectAssetsItcParams struct {
	Ids    []string `validate:"dive,uuid" json:"ids"`
	Cursor string   `json:"cursor"`
	Limit  int      `validate:"gte=0,lte=100" json:"limit"`
}

type InsertAssetItcParams struct {
	Type        assets_dm.Type      `validate:"required,max=32,oneof=CHART INSIGHT AUDIENCE" json:"type"`
	Name        string              `validate:"required,max=128" json:"name"`
	Description string              `validate:"required,max=8192" json:"description"`
	AssetData   assets_dm.AssetData `validate:"required" json:"asset_data"`
}

type UpdateAssetItcParams struct {
	Id          string  `validate:"required,uuid" json:"id"`
	Description *string `validate:"omitempty,max=64" json:"description"`
}

type DeleteAssetItcParams struct {
	Id string `validate:"required,uuid" json:"id"`
}

/// interactor

type AssetsInteractor interface {
	Select(ctx context.Context, params SelectAssetsItcParams) ([]assets_dm.AssetEntity, string, error)
	Insert(ctx context.Context, params ...InsertAssetItcParams) ([]assets_dm.AssetEntity, error)
	Update(ctx context.Context, params ...UpdateAssetItcParams) ([]assets_dm.AssetEntity, error)
	Delete(ctx context.Context, params ...DeleteAssetItcParams) ([]assets_dm.AssetEntity, error)
}

/*
 * Favourites
 */

/// params

type SelectFavouritesItcParams struct {
	Ids      []string `validate:"dive,uuid" json:"ids"`
	UserIds  []string `validate:"dive,uuid" json:"user_ids"`
	AssetIds []string `validate:"dive,uuid" json:"asset_ids"`
	Cursor   string   `json:"cursor"`
	Limit    int      `validate:"gte=0,lte=100" json:"limit"`
}

type InsertFavouriteItcParams struct {
	UserId  string `validate:"required,uuid" json:"user_id"`
	AssetId string `validate:"required,uuid" json:"asset_id"`
}

type DeleteFavouriteItcParams struct {
	Id string `validate:"required,uuid" json:"id"`
}

/// interactor

type FavouritesInteractor interface {
	Select(ctx context.Context, params SelectFavouritesItcParams) ([]favourites_dm.FavouriteEntity, string, error)
	Insert(ctx context.Context, params ...InsertFavouriteItcParams) ([]favourites_dm.FavouriteEntity, error)
	Delete(ctx context.Context, params ...DeleteFavouriteItcParams) ([]favourites_dm.FavouriteEntity, error)
}
