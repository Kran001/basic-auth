package app

import (
	"context"
	"strings"

	"github.com/Kran001/basic-auth/internal/dispatchers"
	"github.com/Kran001/basic-auth/internal/server"
	"github.com/Kran001/basic-auth/internal/store"

	"github.com/Kran001/basic-auth/pkg/config"
	"github.com/Kran001/basic-auth/pkg/database"
	"github.com/Kran001/basic-auth/pkg/logging"
	"github.com/Kran001/basic-auth/pkg/utils"

	"golang.org/x/sync/errgroup"
)

func Run(pathToConfig, sessionsValues string) {
	config, err := config.NewConfig(pathToConfig)
	if err != nil {
		logging.Logger.Fatal("Error while http of path settings configs init. Reason:", err.Error())
	}

	eg, ctx := errgroup.WithContext(context.Background())

	db, err := database.NewDBConnection(config.Database())
	if err != nil {
		logging.Logger.Fatal("Error get init DB connection pool. Reason:", err.Error())
	}

	defer func() {
		if errDBClose := db.Close(); errDBClose != nil {
			logging.Logger.Error("Error closing DB connection pool. Reason:", errDBClose.Error())
		}
	}()

	spoilValues := strings.Split(sessionsValues, "|")
	repositories := store.NewRepositories(db)
	eg.Go(func() error {
		return dispatchers.NewSessionCleaner(repositories.Sessions, spoilValues[0], spoilValues[1]).Run(ctx)
	})

	httpServer := server.NewServer(repositories,
		config,
		config.HTTPSettings().CertLocation(),
		config.HTTPSettings().KeyLocation())

	if err = httpServer.Init(); err != nil {
		logging.Logger.Fatalf("Failed start http server by %s service. Reason:", err.Error())
	}

	eg.Go(func() error {
		if errRun := httpServer.Run(ctx); errRun != nil {
			logging.Logger.Error(errRun.Error())

			return errRun
		}

		return nil
	})

	if err = utils.WaitSignals(ctx); err != nil {
		logging.Logger.Error("Bad finish by wait signal from OS. Reason:", err.Error())
	}

	logging.Logger.Info("Starting shutdown...")
	httpServer.Shutdown(context.Background())
	if err = eg.Wait(); err != nil {
		logging.Logger.Fatal("Goroutine returned with error:", err.Error())
	}

	logging.Logger.Info("Shutdown finished")
}
