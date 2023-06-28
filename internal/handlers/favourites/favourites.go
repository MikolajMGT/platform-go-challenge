package favourites_hl

import (
	favourites_dm "assets/internal/core/domain/favourites"
	"assets/internal/core/ports"
	errs "assets/pkg/errors"
	"assets/pkg/logging"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	webServer     *echo.Echo
	logger        logging.Logger
	favouritesItc ports.FavouritesInteractor
}

func Init(webServer *echo.Echo, logger logging.Logger, interactor ports.FavouritesInteractor) *Handler {

	instance := &Handler{
		webServer:     webServer,
		logger:        logger,
		favouritesItc: interactor,
	}

	instance.webServer.GET("/api/favourites/:id", instance.HandleSelectOne)
	instance.webServer.GET("/api/favourites/user/:userId", instance.HandleSelectMany)
	instance.webServer.POST("/api/favourites/add", instance.HandleInsert)
	instance.webServer.DELETE("/api/favourites/delete/:id", instance.HandleDelete)

	return instance
}

func (h *Handler) HandleSelectOne(ctx echo.Context) (err error) {

	var results []favourites_dm.FavouriteEntity

	h.logger.Info("favourites_hl.HandleSelectOne() performed",
		"id", ctx.Param("id"),
		"results", results,
	)

	id := ctx.Param("id")
	results, _, err = h.favouritesItc.Select(context.Background(), ports.SelectFavouritesItcParams{
		Ids:   []string{id},
		Limit: 1,
	})

	if err == nil && id != "" && len(results) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "record cannot be found")
	} else if err != nil && errors.Is(err, errs.ValidationError) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, results[0])
}

func (h *Handler) HandleSelectMany(ctx echo.Context) (err error) {

	var cursor string
	var nextCursor string
	var results []favourites_dm.FavouriteEntity

	h.logger.Info("favourites_hl.HandleSelectMany() performed",
		"results", results,
	)

	userId := ctx.Param("userId")
	cursor, limit := parseCursorAndLimit(ctx)
	results, nextCursor, err = h.favouritesItc.Select(context.Background(), ports.SelectFavouritesItcParams{
		UserIds: []string{userId},
		Cursor:  cursor,
		Limit:   limit,
	})

	if err != nil && errors.Is(err, errs.ValidationError) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"favourites": results,
		"cursor":     nextCursor,
	})
}

func (h *Handler) HandleInsert(ctx echo.Context) (err error) {
	var results []favourites_dm.FavouriteEntity

	var insertParams ports.InsertFavouriteItcParams
	if err = ctx.Bind(&insertParams); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.logger.Info("favourites_hl.HandleInsert() performed",
		"request", insertParams,
		"results", results,
	)

	results, err = h.favouritesItc.Insert(context.Background(), insertParams)

	if err != nil && errors.Is(err, errs.ValidationError) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil && errors.Is(err, errs.AlreadyExistsError) {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if len(results) == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, results[0])

}

func (h *Handler) HandleDelete(ctx echo.Context) (err error) {

	var results []favourites_dm.FavouriteEntity

	h.logger.Info("favourites_hl.HandleSelectOne() performed",
		"email", ctx.Param("email"),
		"results", results,
	)

	id := ctx.Param("id")
	results, err = h.favouritesItc.Delete(context.Background(), ports.DeleteFavouriteItcParams{
		Id: id,
	})

	if err == nil && id != "" && len(results) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "record cannot be found")
	} else if err != nil && errors.Is(err, errs.ValidationError) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, results[0])
}

func parseCursorAndLimit(ctx echo.Context) (cursor string, limit int) {
	var err error

	cursor = ctx.QueryParam("cursor")
	tmp := ctx.QueryParam("limit")
	if limit, err = strconv.Atoi(tmp); err != nil {
		limit = 0
	}

	return cursor, limit
}
