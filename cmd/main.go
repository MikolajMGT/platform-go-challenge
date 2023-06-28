package main

import (
	"assets/cfg"
	assets_itc "assets/internal/core/interactors/assets"
	favourites_itc "assets/internal/core/interactors/favourites"
	users_itc "assets/internal/core/interactors/users"
	assets_hl "assets/internal/handlers/assets"
	favourites_hl "assets/internal/handlers/favourites"
	users_hl "assets/internal/handlers/users"
	assets_db "assets/internal/repositories/assets"
	audiences_db "assets/internal/repositories/audiences"
	charts_db "assets/internal/repositories/charts"
	favourites_db "assets/internal/repositories/favourites"
	insights_db "assets/internal/repositories/insights"
	sessions_db "assets/internal/repositories/sessions"
	users_db "assets/internal/repositories/users"
	"assets/pkg/logging"
	"assets/pkg/validation"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

const maxRetries = 100

func main() {
	var err error

	/*
	 * Load configs
	 */

	if err = cfg.Configure(cfg.ServiceDefaults()); err != nil {
		panic(errors.Wrap(err, "failed to load configs"))
	}

	/*
	 * Initialize components
	 */

	var webServer *echo.Echo
	if webServer, err = loadComponents(); err != nil {
		panic(errors.Wrap(err, "failed to load components"))
	}

	/*
	 * Serve resources
	 */

	if err := webServer.Start(fmt.Sprintf(":%d", viper.GetInt("webserver.port"))); err != nil {
		panic(fmt.Sprintf("unable to start http server: %s", err))
	}
}

func loadComponents() (webServer *echo.Echo, err error) {
	/// dependencies

	logger := logging.NewDefaultLogger()
	validator := validation.NewDefaultValidator()
	webServer = echo.New()

	webServer.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == http.MethodPost && c.Request().Header.Get("Content-Type") != "application/json" {
				return echo.NewHTTPError(http.StatusBadRequest, "Description-Type must be application/json")
			}
			return next(c)
		}
	})

	var session *gocql.Session
	cluster := gocql.NewCluster(viper.GetString("cassandra.cluster.ip"))
	cluster.Keyspace = "system"

	for i := 0; i < maxRetries; i++ {
		if session, err = cluster.CreateSession(); err != nil {
			log.Println("Could not connect to db: ", err)
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("couldn't connect to database after %v retries", maxRetries))
	}

	createKeySpaceQuery := "CREATE KEYSPACE IF NOT EXISTS %v WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};"
	if err = session.Query(fmt.Sprintf(createKeySpaceQuery, viper.GetString("cassandra.cluster.keyspace"))).Exec(); err != nil {
		return nil, errors.Wrap(err, "failed to create/inspect keyspace")
	}

	session.Close()

	cluster.Keyspace = viper.GetString("cassandra.cluster.keyspace")
	if session, err = cluster.CreateSession(); err != nil {
		panic("failed to connect to database keyspace")
	}

	/// repositories

	usersRepo := users_db.NewCassandraRepo(logger, session)
	chartsRepo := charts_db.NewCassandraRepo(logger, session)
	insightsRepo := insights_db.NewCassandraRepo(logger, session)
	audiencesRepo := audiences_db.NewCassandraRepo(logger, session)
	favouritesRepo := favourites_db.NewCassandraRepo(logger, session)
	assetsRepo := assets_db.NewCassandraRepo(logger, session, chartsRepo, insightsRepo, audiencesRepo)
	sessionsRepo := sessions_db.NewCassandraRepo(logger, session)

	/// interactors
	usersItc := users_itc.NewInteractor(logger, validator, usersRepo)
	favouritesItc := favourites_itc.NewInteractor(logger, validator, favouritesRepo, usersRepo, assetsRepo)
	assetsItc := assets_itc.NewInteractor(logger, validator, assetsRepo, chartsRepo, insightsRepo, audiencesRepo, favouritesRepo)

	/// handlers
	users_hl.Init(webServer, logger, usersItc, sessionsRepo, viper.GetString("auth.secret"))
	favourites_hl.Init(webServer, logger, favouritesItc)
	assets_hl.Init(webServer, logger, assetsItc)

	return webServer, nil
}
