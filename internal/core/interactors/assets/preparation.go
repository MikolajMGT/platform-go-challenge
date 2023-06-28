package assets_itc

import (
	assets_dm "assets/internal/core/domain/assets"
	"assets/internal/core/ports"
	"errors"
)

func prepareCreatableModels(params []ports.InsertAssetItcParams, mapper map[ports.InsertAssetItcParams]string) (results []assets_dm.AssetEntity, err error) {

	for _, param := range params {
		obj := assets_dm.NewAssetEntity()

		obj.Type = param.Type
		obj.Name = param.Name
		obj.Description = param.Description

		if contentId, ok := mapper[param]; ok {
			obj.ContentId = contentId
		} else {
			return nil, errors.New("failed to map data")
		}

		results = append(results, obj)
	}

	return results, nil
}

func prepareUpdatableModels(params []ports.UpdateAssetItcParams, models []assets_dm.AssetEntity) (results []assets_dm.AssetEntity) {
	results = make([]assets_dm.AssetEntity, len(models))
	copy(results, models)

	for idx, param := range params {
		if param.Description != nil && *param.Description != results[idx].Description {
			results[idx].Description = *param.Description
		}
	}

	return results
}
