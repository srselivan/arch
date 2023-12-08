package server

import (
	httpcontroller "arch/internal/server/controller/http"
	natscontroller "arch/internal/server/controller/nats"
	"arch/internal/server/service/auth"
	authrepo "arch/internal/server/service/auth/repo"
	"arch/internal/server/service/message"
	messagerepo "arch/internal/server/service/message/repo"
	"arch/pkg/api/viewerclient"
	"arch/pkg/heartbeat"
	"arch/pkg/logger"
	"arch/pkg/nats"
	"arch/pkg/postgres"
	"arch/pkg/uuid"
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	serviceName   = "server"
	watchDir      = "_heartbeats"
	writeInterval = 5 * time.Second
	fileDir       = "_upload"
)

func Run() {
	log := logger.New("debug", serviceName)

	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		err = os.Mkdir(fileDir, 0777)
		if err != nil {
			log.Fatal().Err(err).Msg("Create directory for upload images error")
		}
	}
	heartbeatService := heartbeat.New(serviceName, watchDir, writeInterval)
	log.Info().Msgf("Run heartbeat for %s directory", watchDir)
	heartbeatService.Run()

	pgConn, err := postgres.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Open postgres connection error")
	}
	if err = postgres.RunMigrations(pgConn.DB); err != nil {
		log.Fatal().Err(err).Msg("Run postgres migrations error")
	}

	natsConn, err := nats.New(nats.Config{
		ClientID: "server_" + uuid.NewV7(),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Create nats connection error")
	}

	viewerClient := viewerclient.New(log, natsConn)
	cacheInstance := cache.New(cache.NoExpiration, cache.NoExpiration)
	authRepo := authrepo.New(pgConn)
	messageRepo := messagerepo.New(cacheInstance)

	messageService := message.New(log, messageRepo, viewerClient)
	authService := auth.New(authRepo, log)

	httpServer := httpcontroller.New(&httpcontroller.Config{
		Addr:           ":8080",
		Log:            log,
		MessageService: messageService,
		AuthService:    authService,
	})

	natsServer := natscontroller.New(natscontroller.Config{
		Conn:           natsConn,
		Log:            log,
		MessageService: messageService,
	})
	go natsServer.Run()

	go func() {
		if err = httpServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Start http server error")
		}
	}()
	log.Info().Str("addr", "8080").Msg("Http server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGKILL)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Shutdown http server error")
	}
	log.Info().Msg("Http server is stopped")

	heartbeatService.Down()
	log.Info().Msg("Heartbeat process is down")

	if err = natsConn.Close(); err != nil {
		log.Fatal().Err(err).Msg("Closing nats connection error")
	}
	log.Info().Msg("Nats connection is closed")

	log.Info().Msg("Server gracefully stopped")
}
