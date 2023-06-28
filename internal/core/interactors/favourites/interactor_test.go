package favourites_itc

import (
	assets_dm "assets/internal/core/domain/assets"
	favourites_dm "assets/internal/core/domain/favourites"
	users_dm "assets/internal/core/domain/users"
	assets_itc "assets/internal/core/interactors/assets"
	users_itc "assets/internal/core/interactors/users"
	"assets/internal/core/ports"
	assets_db "assets/internal/repositories/assets"
	audiences_db "assets/internal/repositories/audiences"
	charts_db "assets/internal/repositories/charts"
	favourites_db "assets/internal/repositories/favourites"
	insights_db "assets/internal/repositories/insights"
	users_db "assets/internal/repositories/users"
	"assets/pkg/logging"
	"assets/pkg/slices"
	"assets/pkg/validation"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InteractorSuite struct {
	suite.Suite
	interactor ports.FavouritesInteractor

	usersItc  ports.UsersInteractor
	assetsItc ports.AssetsInteractor
}

func TestInteractorSuite(t *testing.T) {
	suite.Run(t, new(InteractorSuite))
}

/*
* Tests
 */

/// Select

func (suite *InteractorSuite) TestListShouldReturnErrorWhenInputDataAreIncorrect() {

	params := []ports.SelectFavouritesItcParams{
		// incorrect Ids
		{
			Ids: []string{"fooBar"},
		},
		// incorrect UserIds
		{
			UserIds: []string{"fooBar"},
		},
		// incorrect AssetIds
		{
			AssetIds: []string{"fooBar"},
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

	preparedModels := suite.setupSampleFavourites()
	consideredModels := preparedModels

	type TestCase struct {
		Name    string
		Param   ports.SelectFavouritesItcParams
		Expects []favourites_dm.FavouriteEntity
	}

	testCases := []TestCase{
		{
			Name: "ids",
			Param: ports.SelectFavouritesItcParams{
				Ids: slices.Map(consideredModels, func(model favourites_dm.FavouriteEntity) string { return model.Id }),
			},
			Expects: consideredModels,
		},
		{
			Name: "userIds",
			Param: ports.SelectFavouritesItcParams{
				UserIds: slices.Map(consideredModels, func(model favourites_dm.FavouriteEntity) string { return model.UserId }),
			},
			Expects: consideredModels,
		},
		{
			Name: "assetIds",
			Param: ports.SelectFavouritesItcParams{
				AssetIds: slices.Map(consideredModels, func(model favourites_dm.FavouriteEntity) string { return model.AssetId }),
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

	params := []ports.InsertFavouriteItcParams{
		// missing UserId
		{
			AssetId: uuid.NewString(),
		},
		// missing AssetId
		{
			UserId: uuid.NewString(),
		},
		// incorrect UserId
		{
			UserId:  "fooBar",
			AssetId: uuid.NewString(),
		},
		// incorrect AssetId
		{
			UserId:  uuid.NewString(),
			AssetId: "fooBar",
		},
	}

	for _, param := range params {
		createdModels, err := suite.interactor.Insert(context.Background(), param)

		suite.Empty(createdModels, "should return empty objects list when input data are incorrect")
		suite.ErrorContains(err, "validation error")
	}
}

func (suite *InteractorSuite) TestCreateShouldCreateModelsWithDefaultValues() {

	var (
		users  []users_dm.UserEntity
		assets []assets_dm.AssetEntity
		err    error
	)

	if users, assets, err = suite.setupSampleDependencies(); err != nil {
		panic(err)
	}

	params := []ports.InsertFavouriteItcParams{
		{
			UserId:  users[0].Id,
			AssetId: assets[0].Id,
		},
		{
			UserId:  users[0].Id,
			AssetId: assets[1].Id,
		},
	}

	createdModels, err := suite.interactor.Insert(context.Background(), params...)
	suite.Nil(err, "should return empty error when provided params are correct")

	for _, model := range createdModels {
		suite.AssertValidUuid(model.Id)
		suite.NotZero(model.CreateTime)
		suite.NotZero(model.UpdateTime)
	}

	testsModels, _, err := suite.interactor.Select(context.Background(), ports.SelectFavouritesItcParams{
		Ids: slices.Map(createdModels, func(model favourites_dm.FavouriteEntity) string { return model.Id }),
	})
	suite.ElementsMatch(testsModels, createdModels, "listed and created objects should be the same")
}

/// Delete

func (suite *InteractorSuite) TestDeleteShouldReturnErrorWhenInputDataAreIncorrect() {

	deletedModels, err := suite.interactor.Delete(context.Background(), ports.DeleteFavouriteItcParams{
		Id: "fooBar",
	})

	suite.Empty(deletedModels, "should return empty objects list when input data are incorrect")
	suite.ErrorContains(err, "validation error")
}

func (suite *InteractorSuite) TestDeleteShouldDeleteSpecifiedModel() {

	createdModels := suite.setupSampleFavourites()
	consideredModels := make([]favourites_dm.FavouriteEntity, 2)
	copy(consideredModels, createdModels[:2])

	deletedModels, err := suite.interactor.Delete(context.Background(), []ports.DeleteFavouriteItcParams{
		{Id: consideredModels[0].Id},
		{Id: consideredModels[1].Id},
	}...)
	suite.Nil(err, "should return empty error when provided params are correct")
	suite.EqualValues(consideredModels, deletedModels, "created and deleted objects should be the same")

	models, cursor, err := suite.interactor.Select(context.Background(), ports.SelectFavouritesItcParams{
		Ids: slices.Map(consideredModels, func(model favourites_dm.FavouriteEntity) string { return model.Id }),
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
	favouritesRepo := favourites_db.NewMemoryRepo()
	usersRepo := users_db.NewMemoryRepo()
	assetsRepo := assets_db.NewMemoryRepo()
	chartsRepo := charts_db.NewMemoryRepo()
	insightsRepo := insights_db.NewMemoryRepo()
	audiencesRepo := audiences_db.NewMemoryRepo()

	suite.usersItc = users_itc.NewInteractor(logger, validator, usersRepo)
	suite.assetsItc = assets_itc.NewInteractor(logger, validator, assetsRepo, chartsRepo, insightsRepo, audiencesRepo, favouritesRepo)
	suite.interactor = NewInteractor(logger, validator, favouritesRepo, usersRepo, assetsRepo)
}

func (suite *InteractorSuite) SetupSuite() {
	println("SetupSuite")
}

func (suite *InteractorSuite) SetupTest() {
	println("SetupTest")
	suite.SetupInteractor()
}

func (suite *InteractorSuite) setupSampleFavourites() (models []favourites_dm.FavouriteEntity) {

	var (
		users  []users_dm.UserEntity
		assets []assets_dm.AssetEntity
		err    error
	)

	if users, assets, err = suite.setupSampleDependencies(); err != nil {
		panic(err)
	}

	if models, err = suite.interactor.Insert(context.Background(),
		[]ports.InsertFavouriteItcParams{
			{
				UserId:  users[0].Id,
				AssetId: assets[0].Id,
			},
			{
				UserId:  users[1].Id,
				AssetId: assets[1].Id,
			},
		}...,
	); err != nil {
		panic(err)
	}

	return models
}

func (suite *InteractorSuite) setupSampleDependencies() (users []users_dm.UserEntity, assets []assets_dm.AssetEntity, err error) {

	var user users_dm.UserEntity

	if user, err = suite.usersItc.Register(context.Background(), ports.RegisterUserItcParams{
		Email:    "test123@test.com",
		Password: "test123",
	}); err != nil {
		return nil, nil, err
	}

	users = append(users, user)

	if user, err = suite.usersItc.Register(context.Background(), ports.RegisterUserItcParams{
		Email:    "check123@check123.com",
		Password: "check123",
	}); err != nil {
		return nil, nil, err
	}

	users = append(users, user)

	if assets, err = suite.assetsItc.Insert(context.Background(), []ports.InsertAssetItcParams{
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
	}...); err != nil {
		return nil, nil, err
	}

	return users, assets, err
}

func (suite *InteractorSuite) AssertValidUuid(id string) {
	parsed, err := uuid.Parse(id)
	suite.NotEmpty(parsed)
	suite.Nil(err)
}
