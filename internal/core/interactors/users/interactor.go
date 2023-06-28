package users_itc

import (
	users_dm "assets/internal/core/domain/users"
	"assets/internal/core/ports"
	errs "assets/pkg/errors"
	"assets/pkg/logging"
	"assets/pkg/validation"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Interactor struct {
	logger    logging.Logger
	validator validation.Validator
	usersRepo ports.UsersRepository
}

func NewInteractor(logger logging.Logger, validator validation.Validator, usersRepo ports.UsersRepository) *Interactor {
	return &Interactor{
		logger:    logger,
		validator: validator,
		usersRepo: usersRepo,
	}
}

func (i *Interactor) Login(ctx context.Context, params ports.LoginUserItcParams) (result users_dm.UserEntity, err error) {

	i.logger.Info("users_itc.Login() performed",
		"params", params,
		"result", result,
	)

	if err = i.validator.Validate(params); err != nil {
		return result, errors.Join(errs.ValidationError, err)
	}

	var users []users_dm.UserEntity
	if users, _, err = i.usersRepo.Select(ctx, ports.SelectUsersRepoParams{Emails: []string{params.Email}}); err != nil {
		return result, errors.Join(errs.ProcessingError, err)
	}

	if len(users) < 1 {
		return result, errors.Join(errs.CannotBeFoundError, errors.New("user cannot be found"))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(params.Password)); err != nil {
		return result, errs.AuthenticationError
	}

	return users[0], err
}

func (i *Interactor) Register(ctx context.Context, params ports.RegisterUserItcParams) (result users_dm.UserEntity, err error) {

	i.logger.Info("users_itc.Register() performed",
		"params", params,
		"result", result,
	)

	if err = i.validator.Validate(params); err != nil {
		return result, errors.Join(errs.ValidationError, err)
	}

	var users []users_dm.UserEntity
	if users, _, err = i.usersRepo.Select(ctx, ports.SelectUsersRepoParams{Emails: []string{params.Email}}); err != nil {
		return result, errors.Join(errs.ProcessingError, err)
	}

	if len(users) > 0 {
		return result, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return result, errors.Join(errs.ProcessingError, err)
	}

	obj := users_dm.NewUserEntity()
	obj.Email = params.Email
	obj.Password = string(hashedPassword)

	if _, err = i.usersRepo.Insert(ctx, obj); err != nil {
		return result, errors.Join(errs.ProcessingError, err)
	}

	return obj, err
}
