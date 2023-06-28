package favourites_dm

import (
	"github.com/google/uuid"
	"time"
)

/*
 * Favourite
 */

type Favourite struct {
	UserId  string `validate:"required,uuid" json:"user_id"`
	AssetId string `validate:"required,uuid" json:"asset_id"`
}

type FavouriteEntity struct {
	Favourite
	Id         string    `validate:"required,uuid" json:"id"`
	CreateTime time.Time `validate:"required" json:"create_time"`
	UpdateTime time.Time `validate:"required" json:"update_time"`
}

func NewFavouriteEntity() FavouriteEntity {
	now := time.Now()

	return FavouriteEntity{
		Id:         uuid.NewString(),
		CreateTime: now,
		UpdateTime: now,
	}
}
