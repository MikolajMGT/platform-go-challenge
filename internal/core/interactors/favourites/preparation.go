package favourites_itc

import (
	favourites_dm "assets/internal/core/domain/favourites"
	"assets/internal/core/ports"
)

func prepareCreatableModels(params []ports.InsertFavouriteItcParams) (results []favourites_dm.FavouriteEntity) {

	for _, param := range params {
		obj := favourites_dm.NewFavouriteEntity()

		obj.UserId = param.UserId
		obj.AssetId = param.AssetId

		results = append(results, obj)
	}

	return results
}
