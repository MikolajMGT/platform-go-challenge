package assets_itc

import (
	assets_dm "assets/internal/core/domain/assets"
	"assets/internal/core/ports"
	assets_db "assets/internal/repositories/assets"
	audiences_db "assets/internal/repositories/audiences"
	charts_db "assets/internal/repositories/charts"
	favourites_db "assets/internal/repositories/favourites"
	insights_db "assets/internal/repositories/insights"
	"assets/pkg/logging"
	"assets/pkg/slices"
	"assets/pkg/validation"
	"context"
	r "crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
	"time"
)

type InteractorSuite struct {
	suite.Suite
	interactor ports.AssetsInteractor
}

func TestInteractorSuite(t *testing.T) {
	suite.Run(t, new(InteractorSuite))
}

/*
* Tests
 */

/// Select

func (suite *InteractorSuite) TestListShouldReturnErrorWhenInputDataAreIncorrect() {

	params := []ports.SelectAssetsItcParams{
		// incorrect Ids
		{
			Ids: []string{"fooBar"},
		},
		// incorrect Limit
		{
			Limit: -1,
		},
	}

	for _, param := range params {
		testsModels, cursor, err := suite.interactor.Select(context.Background(), param)

		suite.Empty(testsModels, "should return empty objects list when params are empty")
		suite.Equal("", cursor, "should empty cursor when nothing found")
		suite.ErrorContains(err, "validation error")
	}
}

func (suite *InteractorSuite) TestListShouldReturnModelListByFilters() {

	preparedModels := suite.setupSampleAssets()
	consideredModels := preparedModels[:2]

	type TestCase struct {
		Name    string
		Param   ports.SelectAssetsItcParams
		Expects []assets_dm.AssetEntity
	}

	testCases := []TestCase{
		{
			Name: "ids",
			Param: ports.SelectAssetsItcParams{
				Ids: slices.Map(consideredModels, func(model assets_dm.AssetEntity) string { return model.Id }),
			},
			Expects: consideredModels,
		},
	}

	for _, c := range testCases {
		testsModels, cursor, err := suite.interactor.Select(context.Background(), c.Param)

		suite.Nil(err, fmt.Sprintf("should return empty error when %s was specified", c.Name))
		suite.Equal("", cursor, "should return empty cursor when all data was returned")
		suite.ElementsMatch(testsModels, c.Expects, fmt.Sprintf("should return only objects list with specfied %s", c.Name))
	}
}

/// Insert

func (suite *InteractorSuite) TestCreateShouldReturnErrorWhenInputDataAreIncorrect() {

	params := []ports.InsertAssetItcParams{
		// missing Type
		{
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Insight: &assets_dm.Insight{
					Text: "Nice Insight",
				},
			},
		},
		// missing Name
		{
			Type:        assets_dm.TypeInsight,
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Insight: &assets_dm.Insight{
					Text: "Nice Insight",
				},
			},
		},
		// missing Description
		{
			Type: assets_dm.TypeInsight,
			Name: "Nice Name",
			AssetData: assets_dm.AssetData{
				Insight: &assets_dm.Insight{
					Text: "Nice Insight",
				},
			},
		},
		// missing AssetData
		{
			Type:        assets_dm.TypeInsight,
			Name:        "Nice Name",
			Description: "Nice Description",
		},
		// incorrect Type
		{
			Type:        "FooBar",
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Insight: &assets_dm.Insight{
					Text: "Nice Insight",
				},
			},
		},
		// incorrect Name
		{
			Type:        assets_dm.TypeInsight,
			Name:        "Nice Name Nice Name Nice Name Nice Name Nice Name Nice Name Nice Name Nice Name Nice Name Nice Name Nice Name Nice Name Nice Name Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Insight: &assets_dm.Insight{
					Text: "Nice Insight",
				},
			},
		},
		// incorrect Description
		{
			Type:        assets_dm.TypeInsight,
			Name:        "Nice Name",
			Description: randomString(9000),
			AssetData: assets_dm.AssetData{
				Insight: &assets_dm.Insight{
					Text: "Nice Insight",
				},
			},
		},
		// incorrect Insight Text
		{
			Type:        assets_dm.TypeInsight,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Insight: &assets_dm.Insight{
					Text: randomString(5000),
				},
			},
		},
		// missing Chart ChartTitle
		{
			Type:        assets_dm.TypeChart,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Chart: &assets_dm.Chart{
					XAxisTitle: "important axis",
					YAxisTitle: "second axis",
					Data:       "something",
				},
			},
		},
		// missing Chart XAxisTitle
		{
			Type:        assets_dm.TypeChart,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Chart: &assets_dm.Chart{
					ChartTitle: "nice title",
					YAxisTitle: "second axis",
					Data:       "something",
				},
			},
		},
		// missing Chart YAxisTitle
		{
			Type:        assets_dm.TypeChart,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Chart: &assets_dm.Chart{
					ChartTitle: "nice title",
					XAxisTitle: "important axis",
					Data:       "something",
				},
			},
		},
		// missing Chart Data
		{
			Type:        assets_dm.TypeChart,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Chart: &assets_dm.Chart{
					ChartTitle: "nice title",
					XAxisTitle: "important axis",
					YAxisTitle: "second axis",
				},
			},
		},
		// incorrect Chart Title
		{
			Type:        assets_dm.TypeChart,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Chart: &assets_dm.Chart{
					ChartTitle: randomString(200),
					XAxisTitle: "important axis",
					YAxisTitle: "second axis",
					Data:       "something",
				},
			},
		},
		// incorrect XAxisTitle
		{
			Type:        assets_dm.TypeChart,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Chart: &assets_dm.Chart{
					ChartTitle: "nice title",
					XAxisTitle: randomString(200),
					YAxisTitle: "second axis",
					Data:       "something",
				},
			},
		},
		// incorrect YAxisTitle
		{
			Type:        assets_dm.TypeChart,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Chart: &assets_dm.Chart{
					ChartTitle: "nice title",
					XAxisTitle: "important axis",
					YAxisTitle: randomString(200),
					Data:       "something",
				},
			},
		},
		// missing Audience Gender
		{
			Type:        assets_dm.TypeAudience,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Audience: &assets_dm.Audience{
					BirthCountry:       "Greece",
					AgeGroup:           "18-23",
					SocialMediaHours:   3,
					PurchasesLastMonth: 6,
				},
			},
		},
		// missing Audience BirthCountry
		{
			Type:        assets_dm.TypeAudience,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Audience: &assets_dm.Audience{
					Gender:             assets_dm.GenderFemale,
					AgeGroup:           "18-23",
					SocialMediaHours:   3,
					PurchasesLastMonth: 6,
				},
			},
		},
		// missing Audience AgeGroup
		{
			Type:        assets_dm.TypeAudience,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Audience: &assets_dm.Audience{
					Gender:             assets_dm.GenderFemale,
					BirthCountry:       "Greece",
					SocialMediaHours:   3,
					PurchasesLastMonth: 6,
				},
			},
		},
		// incorrect Audience Gender
		{
			Type:        assets_dm.TypeAudience,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Audience: &assets_dm.Audience{
					Gender:             "Foobar",
					BirthCountry:       "Greece",
					AgeGroup:           "18-23",
					SocialMediaHours:   3,
					PurchasesLastMonth: 6,
				},
			},
		},
		// incorrect Audience BirthCountry
		{
			Type:        assets_dm.TypeAudience,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Audience: &assets_dm.Audience{
					Gender:             assets_dm.GenderFemale,
					BirthCountry:       randomString(200),
					AgeGroup:           "18-23",
					SocialMediaHours:   3,
					PurchasesLastMonth: 6,
				},
			},
		},
		// incorrect Audience AgeGroup
		{
			Type:        assets_dm.TypeAudience,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Audience: &assets_dm.Audience{
					Gender:             assets_dm.GenderFemale,
					BirthCountry:       "Greece",
					AgeGroup:           "FooBar",
					SocialMediaHours:   3,
					PurchasesLastMonth: 6,
				},
			},
		},
		// incorrect Social Media Hours
		{
			Type:        assets_dm.TypeAudience,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Audience: &assets_dm.Audience{
					Gender:             "Foobar",
					BirthCountry:       "Greece",
					AgeGroup:           "18-23",
					SocialMediaHours:   -1,
					PurchasesLastMonth: 6,
				},
			},
		},
		// incorrect Audience Purchases Last Month
		{
			Type:        assets_dm.TypeAudience,
			Name:        "Nice Name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Audience: &assets_dm.Audience{
					Gender:             "Foobar",
					BirthCountry:       "Greece",
					AgeGroup:           "18-23",
					SocialMediaHours:   3,
					PurchasesLastMonth: -1,
				},
			},
		},
	}

	for _, param := range params {
		createdModels, err := suite.interactor.Insert(context.Background(), param)

		suite.Empty(createdModels, "should return empty objects list when input data are incorrect")
		suite.ErrorContains(err, "validation error")
	}
}

func (suite *InteractorSuite) TestCreateShouldCreateModelsWithDefaultValues() {

	params := []ports.InsertAssetItcParams{
		{
			Type:        assets_dm.TypeInsight,
			Name:        "test name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Insight: &assets_dm.Insight{
					Text: "Nice Insight",
				},
			},
		},
		{
			Type:        assets_dm.TypeChart,
			Name:        "test name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Chart: &assets_dm.Chart{
					ChartTitle: "interesting title",
					XAxisTitle: "important parameter",
					YAxisTitle: "also important thing",
					Data:       "test",
				},
			},
		},
		{
			Type:        assets_dm.TypeAudience,
			Name:        "test name",
			Description: "Nice Description",
			AssetData: assets_dm.AssetData{
				Audience: &assets_dm.Audience{
					Gender:             assets_dm.GenderFemale,
					BirthCountry:       "Greece",
					AgeGroup:           "18-23",
					SocialMediaHours:   3,
					PurchasesLastMonth: 6,
				},
			},
		},
	}

	createdModels, err := suite.interactor.Insert(context.Background(), params...)
	suite.Nil(err, "should return empty error when provided params are correct")

	for _, model := range createdModels {
		suite.AssertValidUuid(model.Id)
		suite.NotZero(model.CreateTime)
		suite.NotZero(model.UpdateTime)
	}

	testsModels, _, err := suite.interactor.Select(context.Background(), ports.SelectAssetsItcParams{
		Ids: slices.Map(createdModels, func(model assets_dm.AssetEntity) string { return model.Id }),
	})
	suite.ElementsMatch(testsModels, createdModels, "listed and created objects should be the same")
}

/// Update

func (suite *InteractorSuite) TestUpdateShouldReturnErrorWhenInputDataAreIncorrect() {

	preparedModels := suite.setupSampleAssets()
	incorrectDescription := randomString(9000)
	params := []ports.UpdateAssetItcParams{
		// incorrect Id
		{
			Id: "fooBar",
		},
		// incorrect Description
		{
			Id:          preparedModels[0].Id,
			Description: &incorrectDescription,
		},
	}

	for _, param := range params {
		updatedModels, err := suite.interactor.Update(context.Background(), param)

		suite.Empty(updatedModels, "should return empty objects list when input data are incorrect")
		suite.ErrorContains(err, "validation error")
	}
}

func (suite *InteractorSuite) TestUpdateShouldUpdateSpecifiedModel() {

	createdModels := suite.setupSampleAssets()
	newDescription1 := "New Description 1"
	newDescription2 := "New Description 2"
	newDescription3 := "New Description 3"

	params := []ports.UpdateAssetItcParams{
		{
			Id:          createdModels[0].Id,
			Description: &newDescription1,
		},
		{
			Id:          createdModels[1].Id,
			Description: &newDescription2,
		},
		{
			Id:          createdModels[2].Id,
			Description: &newDescription3,
		},
	}

	updatedModels, err := suite.interactor.Update(context.Background(), params...)
	expectedModels := make([]assets_dm.AssetEntity, len(updatedModels))
	copy(expectedModels, updatedModels)

	for k := range expectedModels {

		if params[k].Description != nil {
			expectedModels[k].Description = *params[k].Description
		}
	}
	suite.Nil(err, "should return empty error when provided params are correct")
	suite.Equal(expectedModels, updatedModels, "expected and returned objects should be the same")

	testsModels, cursor, err := suite.interactor.Select(context.Background(), ports.SelectAssetsItcParams{
		Ids: slices.Map(expectedModels, func(model assets_dm.AssetEntity) string { return model.Id }),
	})
	suite.Empty(err, "it shouldn't be an error while listing models")
	suite.Empty(cursor, "cursor should be empty when all data was fetched")
	suite.ElementsMatch(testsModels, updatedModels, "listed and updated objects should be the same")
}

/// Delete

func (suite *InteractorSuite) TestDeleteShouldReturnErrorWhenInputDataAreIncorrect() {

	deletedModels, err := suite.interactor.Delete(context.Background(), ports.DeleteAssetItcParams{
		Id: "fooBar",
	})

	suite.Empty(deletedModels, "should return empty objects list when input data are incorrect")
	suite.ErrorContains(err, "validation error")
}

func (suite *InteractorSuite) TestDeleteShouldDeleteSpecifiedModel() {

	createdModels := suite.setupSampleAssets()
	consideredModels := make([]assets_dm.AssetEntity, 2)
	copy(consideredModels, createdModels[:2])

	deletedModels, err := suite.interactor.Delete(context.Background(), []ports.DeleteAssetItcParams{
		{Id: consideredModels[0].Id},
		{Id: consideredModels[1].Id},
	}...)
	suite.Nil(err, "should return empty error when provided params are correct")
	suite.EqualValues(consideredModels, deletedModels, "created and deleted objects should be the same")

	models, cursor, err := suite.interactor.Select(context.Background(), ports.SelectAssetsItcParams{
		Ids: slices.Map(consideredModels, func(model assets_dm.AssetEntity) string { return model.Id }),
	})
	suite.Nil(err, "error should be nil")
	suite.Empty(cursor, "cursor should be empty when all data are fetched")
	suite.Equal(0, len(models), "count number after deletion should be equal zero")
}

/*
* SUITE SETUP
 */

func (suite *InteractorSuite) SetupInteractor() {
	logger := logging.NewDefaultLogger()
	validator := validation.NewDefaultValidator()
	assetsRepo := assets_db.NewMemoryRepo()
	chartsRepo := charts_db.NewMemoryRepo()
	insightsRepo := insights_db.NewMemoryRepo()
	audiencesRepo := audiences_db.NewMemoryRepo()
	favouritesRepo := favourites_db.NewMemoryRepo()

	suite.interactor = NewInteractor(logger, validator, assetsRepo, chartsRepo, insightsRepo, audiencesRepo, favouritesRepo)
}

func (suite *InteractorSuite) SetupSuite() {
	println("SetupSuite")
}

func (suite *InteractorSuite) SetupTest() {
	println("SetupTest")
	suite.SetupInteractor()
}

func (suite *InteractorSuite) setupSampleAssets() (models []assets_dm.AssetEntity) {

	var err error
	if models, err = suite.interactor.Insert(context.Background(),
		[]ports.InsertAssetItcParams{
			{
				Type:        assets_dm.TypeInsight,
				Name:        "test name",
				Description: "Nice Description",
				AssetData: assets_dm.AssetData{
					Insight: &assets_dm.Insight{
						Text: "Nice Insight",
					},
				},
			},
			{
				Type:        assets_dm.TypeChart,
				Name:        "test name",
				Description: "Nice Description",
				AssetData: assets_dm.AssetData{
					Chart: &assets_dm.Chart{
						ChartTitle: "interesting title",
						XAxisTitle: "important parameter",
						YAxisTitle: "also important thing",
						Data:       "test",
					},
				},
			},
			{
				Type:        assets_dm.TypeAudience,
				Name:        "test name",
				Description: "Nice Description",
				AssetData: assets_dm.AssetData{
					Audience: &assets_dm.Audience{
						Gender:             assets_dm.GenderFemale,
						BirthCountry:       "Greece",
						AgeGroup:           "18-23",
						SocialMediaHours:   3,
						PurchasesLastMonth: 6,
					},
				},
			},
		}...,
	); err != nil {
		panic(err)
	}

	return models
}

func (suite *InteractorSuite) AssertValidUuid(id string) {
	parsed, err := uuid.Parse(id)
	suite.NotEmpty(parsed)
	suite.Nil(err)
}

func randomString(length int) string {
	rand.NewSource(time.Now().UnixNano())
	b := make([]byte, length+2)
	_, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)[2 : length+2]
}
