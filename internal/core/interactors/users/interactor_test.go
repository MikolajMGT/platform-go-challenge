package users_itc

import (
	users_dm "assets/internal/core/domain/users"
	"assets/internal/core/ports"
	users_db "assets/internal/repositories/users"
	"assets/pkg/logging"
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
	interactor ports.UsersInteractor
}

func TestInteractorSuite(t *testing.T) {
	suite.Run(t, new(InteractorSuite))
}

/*
* Tests
 */

/// Login

func (suite *InteractorSuite) TestLoginShouldReturnErrorWhenInputDataAreIncorrect() {

	params := []ports.LoginUserItcParams{
		{
			Email:    "fooBar",
			Password: "topSecret",
		},
		{
			Email:    "test@test.test",
			Password: randomString(1000),
		},
	}

	for _, param := range params {
		testsModels, err := suite.interactor.Login(context.Background(), param)

		suite.Empty(testsModels, "should return empty objects list when params are empty")
		suite.ErrorContains(err, "validation error")
	}
}

func (suite *InteractorSuite) TestLoginShouldReturnUserObject() {

	preparedModel := suite.setupSampleUser()

	testModel, err := suite.interactor.Login(context.Background(), ports.LoginUserItcParams{
		Email:    preparedModel.Email,
		Password: "test123",
	})

	suite.Nil(err, fmt.Sprintf("should return empty error "))
	suite.Equal(preparedModel, testModel, "objects should match")
}

/// Register

func (suite *InteractorSuite) TestRegisterShouldReturnErrorWhenInputDataAreIncorrect() {

	params := []ports.RegisterUserItcParams{
		// missing Email
		{
			Password: "test123",
		},
		// missing Password
		{
			Email: "test@test.com",
		},
		// incorrect Email
		{
			Email:    "fooBar",
			Password: "test123",
		},
		// incorrect Password
		{
			Email:    "test@test.com",
			Password: randomString(900),
		},
	}

	for _, param := range params {
		createdModels, err := suite.interactor.Register(context.Background(), param)

		suite.Empty(createdModels, "should return empty objects list when input data are incorrect")
		suite.ErrorContains(err, "validation error")
	}
}

func (suite *InteractorSuite) TestRegisterShouldReturnObject() {

	testModel, err := suite.interactor.Register(context.Background(), ports.RegisterUserItcParams{
		Email:    "test@test.com",
		Password: "test123",
	})

	suite.Nil(err, fmt.Sprintf("should return empty error "))
	suite.NotEmpty(testModel.Id, "field should be populated")
	suite.Equal(testModel.Email, testModel.Email, "field should match")
	suite.Equal(testModel.Password, testModel.Password, "field should match")
	suite.NotEmpty(testModel.CreateTime, "field should be populated")
	suite.NotEmpty(testModel.UpdateTime, "field should be populated")
}

/*
* SUITE SETUP
 */

func (suite *InteractorSuite) SetupInteractor() {
	logger := logging.NewDefaultLogger()
	validator := validation.NewDefaultValidator()
	usersRepo := users_db.NewMemoryRepo()

	suite.interactor = NewInteractor(logger, validator, usersRepo)
}

func (suite *InteractorSuite) SetupSuite() {
	println("SetupSuite")
}

func (suite *InteractorSuite) SetupTest() {
	println("SetupTest")
	suite.SetupInteractor()
}

func (suite *InteractorSuite) setupSampleUser() (model users_dm.UserEntity) {

	var err error
	if model, err = suite.interactor.Register(context.Background(), ports.RegisterUserItcParams{
		Email:    "test@test.com",
		Password: "test123",
	}); err != nil {
		panic(err)
	}

	return model
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
