package users_hl

import (
	users_dm "assets/internal/core/domain/users"
	"assets/internal/core/ports"
	errs "assets/pkg/errors"
	"assets/pkg/logging"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Handler struct {
	webServer   *echo.Echo
	logger      logging.Logger
	usersItc    ports.UsersInteractor
	cookieStore sessions.Store
	secret      string
}

func Init(webServer *echo.Echo, logger logging.Logger, interactor ports.UsersInteractor, cookieStore sessions.Store, secret string) *Handler {

	instance := &Handler{
		webServer:   webServer,
		logger:      logger,
		usersItc:    interactor,
		cookieStore: cookieStore,
		secret:      secret,
	}

	instance.webServer.POST("/api/users/login", instance.HandleLogin)
	instance.webServer.POST("/api/users/register", instance.HandleRegister)

	return instance
}

func (h *Handler) HandleLogin(ctx echo.Context) (err error) {
	var result users_dm.UserEntity

	var loginParams ports.LoginUserItcParams
	if err = ctx.Bind(&loginParams); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.logger.Info("users_hl.HandleLogin() performed",
		"request", loginParams,
		"result", result,
	)

	result, err = h.usersItc.Login(context.Background(), loginParams)

	if err != nil && errors.Is(err, errs.ValidationError) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil && errors.Is(err, errs.AuthenticationError) {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = h.authenticate(ctx, result.Id, result.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	result.Password = "<top secret>"

	return ctx.JSON(http.StatusOK, result)
}

func (h *Handler) HandleRegister(ctx echo.Context) (err error) {
	var result users_dm.UserEntity

	var insertParams ports.RegisterUserItcParams
	if err = ctx.Bind(&insertParams); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.logger.Info("users_hl.HandleRegister() performed",
		"request", insertParams,
		"result", result,
	)

	result, err = h.usersItc.Register(context.Background(), insertParams)

	if err != nil && errors.Is(err, errs.ValidationError) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if errors.Is(err, errors.New("user already exists")) {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = h.authenticate(ctx, result.Id, result.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	result.Password = "<top secret>"

	return ctx.JSON(http.StatusCreated, result)

}

func (h *Handler) authenticate(ctx echo.Context, userId, email string) (err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.secret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	session, err := h.cookieStore.Get(ctx.Request(), "session-name")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	session.Values["token"] = tokenString
	if err = session.Save(ctx.Request(), ctx.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
