package assets_hl

import (
	assets_dm "assets/internal/core/domain/assets"
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
	webServer *echo.Echo
	logger    logging.Logger
	assetsItc ports.AssetsInteractor
}

func Init(webServer *echo.Echo, logger logging.Logger, interactor ports.AssetsInteractor) *Handler {

	instance := &Handler{
		webServer: webServer,
		logger:    logger,
		assetsItc: interactor,
	}

	instance.webServer.GET("/api/assets", instance.HandleSelectMany)
	instance.webServer.GET("/api/assets/:id", instance.HandleSelectOne)
	instance.webServer.POST("/api/assets/create", instance.HandleInsert)
	instance.webServer.PATCH("/api/assets/update", instance.HandleUpdate)
	instance.webServer.DELETE("/api/assets/delete/:id", instance.HandleDelete)

	return instance
}

func (h *Handler) HandleSelectOne(ctx echo.Context) (err error) {

	var results []assets_dm.AssetEntity

	h.logger.Info("assets_hl.HandleSelectOne() performed",
		"id", ctx.Param("id"),
		"results", results,
	)

	id := ctx.Param("id")
	results, _, err = h.assetsItc.Select(context.Background(), ports.SelectAssetsItcParams{
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
	var results []assets_dm.AssetEntity

	h.logger.Info("assets_hl.HandleSelectMany() performed",
		"results", results,
	)

	cursor, limit := parseCursorAndLimit(ctx)
	results, nextCursor, err = h.assetsItc.Select(context.Background(), ports.SelectAssetsItcParams{
		Cursor: cursor,
		Limit:  limit,
	})

	if err != nil && errors.Is(err, errs.ValidationError) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"assets": results,
		"cursor": nextCursor,
	})
}

func (h *Handler) HandleInsert(ctx echo.Context) (err error) {
	var results []assets_dm.AssetEntity

	var insertParams ports.InsertAssetItcParams
	if err = ctx.Bind(&insertParams); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.logger.Info("assets_hl.HandleInsert() performed",
		"request", insertParams,
		"results", results,
	)

	results, err = h.assetsItc.Insert(context.Background(), insertParams)

	if err != nil && errors.Is(err, errs.ValidationError) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if len(results) == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, results[0])

}

func (h *Handler) HandleUpdate(ctx echo.Context) (err error) {
	var results []assets_dm.AssetEntity

	var updateParams ports.UpdateAssetItcParams
	if err = ctx.Bind(&updateParams); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.logger.Info("assets_hl.HandleUpdate() performed",
		"request", updateParams,
		"results", results,
	)

	results, err = h.assetsItc.Update(context.Background(), updateParams)

	if len(results) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "record cannot be found")
	} else if err != nil && errors.Is(err, errs.ValidationError) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, results[0])

}

func (h *Handler) HandleDelete(ctx echo.Context) (err error) {

	var results []assets_dm.AssetEntity

	h.logger.Info("assets_hl.HandleSelectOne() performed",
		"email", ctx.Param("email"),
		"results", results,
	)

	id := ctx.Param("id")
	results, err = h.assetsItc.Delete(context.Background(), ports.DeleteAssetItcParams{
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
