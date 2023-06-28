package favourites_itc

import "assets/internal/core/ports"

func convertSelectParams(params ports.SelectFavouritesItcParams) (result ports.SelectFavouritesRepoParams) {
	return ports.SelectFavouritesRepoParams{
		Ids:      params.Ids,
		UserIds:  params.UserIds,
		AssetIds: params.AssetIds,
		Cursor:   params.Cursor,
		Limit:    params.Limit,
	}
}
