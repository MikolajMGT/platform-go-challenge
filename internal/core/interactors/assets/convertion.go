package assets_itc

import "assets/internal/core/ports"

func convertSelectParams(params ports.SelectAssetsItcParams) (result ports.SelectAssetsRepoParams) {
	return ports.SelectAssetsRepoParams{
		Ids:    params.Ids,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	}
}
