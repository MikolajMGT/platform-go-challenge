package assets_dm

import (
	"github.com/google/uuid"
	"time"
)

/*
 * AssetData
 */

type Asset struct {
	ContentId   string            `validate:"required,uuid" json:"content_id"`
	Type        Type              `validate:"required,max=32,oneof=CHART INSIGHT AUDIENCE" json:"type"`
	Name        string            `validate:"required,max=32" json:"name"`
	Description string            `validate:"required,max=8192" json:"description"`
	AssetData   AssetDataEntities `validate:"required" json:"asset_data"`
}

type AssetEntity struct {
	Asset
	Id         string    `validate:"required,uuid" json:"id"`
	CreateTime time.Time `validate:"required" json:"create_time"`
	UpdateTime time.Time `validate:"required" json:"update_time"`
}

type AssetData struct {
	Chart    *Chart    `json:"chart,omitempty"`
	Insight  *Insight  `json:"insight,omitempty"`
	Audience *Audience `json:"audience,omitempty"`
}

type AssetDataEntities struct {
	Chart    *ChartEntity    `json:"chart,omitempty"`
	Insight  *InsightEntity  `json:"insight,omitempty"`
	Audience *AudienceEntity `json:"audience,omitempty"`
}

func NewAssetEntity() AssetEntity {
	now := time.Now()

	return AssetEntity{
		Id:         uuid.NewString(),
		CreateTime: now,
		UpdateTime: now,
	}
}

/*
 * ChartEntity
 */

type Chart struct {
	ChartTitle string `validate:"required,max=32" json:"title"`
	XAxisTitle string `validate:"required,max=32" json:"x_axis_title"`
	YAxisTitle string `validate:"required,max=32" json:"y_axis_title"`
	Data       any    `validate:"required" json:"data"`
}

type ChartEntity struct {
	Chart
	Id         string    `validate:"required,uuid" json:"id"`
	CreateTime time.Time `validate:"required" json:"create_time"`
	UpdateTime time.Time `validate:"required" json:"update_time"`
}

func NewChartEntity() ChartEntity {
	now := time.Now()

	return ChartEntity{
		Id:         uuid.NewString(),
		CreateTime: now,
		UpdateTime: now,
	}
}

/*
 * InsightEntity
 */

type Insight struct {
	Text string `validate:"required,max=1024" json:"text"`
}

type InsightEntity struct {
	Insight
	Id         string    `validate:"required,uuid" json:"id"`
	CreateTime time.Time `validate:"required" json:"create_time"`
	UpdateTime time.Time `validate:"required" json:"update_time"`
}

func NewInsightEntity() InsightEntity {
	now := time.Now()

	return InsightEntity{
		Id:         uuid.NewString(),
		CreateTime: now,
		UpdateTime: now,
	}
}

/*
 * AudienceEntity
 */

type Audience struct {
	Gender             Gender   `validate:"required,max=32,oneof=MALE FEMALE" json:"gender"`
	BirthCountry       string   `validate:"required,max=32" json:"birth_country"`
	AgeGroup           AgeGroup `validate:"required,max=32,oneof=18-23 24-35 36-45 46+" json:"age_group"`
	SocialMediaHours   int64    `validate:"gt=-1" json:"social_media_hours"`
	PurchasesLastMonth int64    `validate:"gt=-1" json:"purchases_last_month"`
}

type AudienceEntity struct {
	Audience
	Id         string    `validate:"required,uuid" json:"id"`
	CreateTime time.Time `validate:"required" json:"create_time"`
	UpdateTime time.Time `validate:"required" json:"update_time"`
}

func NewAudienceEntity() AudienceEntity {
	now := time.Now()

	return AudienceEntity{
		Id:         uuid.NewString(),
		CreateTime: now,
		UpdateTime: now,
	}
}
