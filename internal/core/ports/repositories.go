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

type SelectUsersRepoParams struct {
	Ids    []string `validate:"dive,uuid" json:"ids"`
	Emails []string `validate:"dive,email" json:"emails"`
	Cursor string   `json:"cursor"`
	Limit  int      `validate:"gte=0,lte=100" json:"limit"`
}

/// repository

type UsersRepository interface {
	Select(ctx context.Context, params SelectUsersRepoParams) ([]users_dm.UserEntity, string, error)
	Insert(ctx context.Context, models ...users_dm.UserEntity) ([]users_dm.UserEntity, error)
}

/*
 * Assets
 */

/// params

type SelectAssetsRepoParams struct {
	Ids    []string
	Cursor string
	Limit  int
}

/// repository

type AssetsRepository interface {
	Select(ctx context.Context, params SelectAssetsRepoParams) ([]assets_dm.AssetEntity, string, error)
	Insert(ctx context.Context, models ...assets_dm.AssetEntity) ([]assets_dm.AssetEntity, error)
	Update(ctx context.Context, models ...assets_dm.AssetEntity) ([]assets_dm.AssetEntity, error)
	Delete(ctx context.Context, models ...assets_dm.AssetEntity) ([]assets_dm.AssetEntity, error)
}

/*
 * Charts
 */

/// params

type SelectChartsRepoParams struct {
	Ids    []string
	Cursor string
	Limit  int
}

/// repository

type ChartsRepository interface {
	Select(ctx context.Context, params SelectChartsRepoParams) ([]assets_dm.ChartEntity, string, error)
	Insert(ctx context.Context, models ...assets_dm.ChartEntity) ([]assets_dm.ChartEntity, error)
	Delete(ctx context.Context, models ...assets_dm.ChartEntity) ([]assets_dm.ChartEntity, error)
}

/*
 * Insights
 */

/// params

type SelectInsightsRepoParams struct {
	Ids    []string
	Cursor string
	Limit  int
}

/// repository

type InsightsRepository interface {
	Select(ctx context.Context, params SelectInsightsRepoParams) ([]assets_dm.InsightEntity, string, error)
	Insert(ctx context.Context, models ...assets_dm.InsightEntity) ([]assets_dm.InsightEntity, error)
	Delete(ctx context.Context, models ...assets_dm.InsightEntity) ([]assets_dm.InsightEntity, error)
}

/*
 * Audiences
 */

/// params

type SelectAudiencesRepoParams struct {
	Ids    []string
	Cursor string
	Limit  int
}

/// repository

type AudiencesRepository interface {
	Select(ctx context.Context, params SelectAudiencesRepoParams) ([]assets_dm.AudienceEntity, string, error)
	Insert(ctx context.Context, models ...assets_dm.AudienceEntity) ([]assets_dm.AudienceEntity, error)
	Delete(ctx context.Context, models ...assets_dm.AudienceEntity) ([]assets_dm.AudienceEntity, error)
}

/*
 * Favourites
 */

/// params

type SelectFavouritesRepoParams struct {
	Ids      []string
	UserIds  []string
	AssetIds []string
	Cursor   string
	Limit    int
}

/// repository

type FavouritesRepository interface {
	Select(ctx context.Context, params SelectFavouritesRepoParams) ([]favourites_dm.FavouriteEntity, string, error)
	Insert(ctx context.Context, models ...favourites_dm.FavouriteEntity) ([]favourites_dm.FavouriteEntity, error)
	Delete(ctx context.Context, models ...favourites_dm.FavouriteEntity) ([]favourites_dm.FavouriteEntity, error)
}
