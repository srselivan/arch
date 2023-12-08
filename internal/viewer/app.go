package viewer

import (
	natscontroller "arch/internal/viewer/controller/nats"
	"arch/internal/viewer/service/message"
	"arch/internal/viewer/service/message/repo"
	"arch/internal/viewer/service/screenupdater"
	"arch/pkg/api/consolepresenter"
	"arch/pkg/api/serverclient"
	"arch/pkg/logger"
	"arch/pkg/nats"
	"arch/pkg/uuid"
	"github.com/patrickmn/go-cache"
	"os"
	"os/signal"
	"syscall"
)

const (
	serviceName = "viewer"
)

func Run() {
	log := logger.New("debug", serviceName)

	natsConn, err := nats.New(nats.Config{
		ClientID: "viewer_" + uuid.NewV7(),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Create nats connection error")
	}

	serverClient := serverclient.New(log, natsConn)
	consolePresenter := consolepresenter.New()

	cacheInstance := cache.New(cache.NoExpiration, cache.NoExpiration)
	messageRepo := repo.New(cacheInstance)

	messageService := message.New(log, messageRepo)
	consoleUpdater := screenupdater.New(log, messageRepo, serverClient, consolePresenter)

	natsServer := natscontroller.New(natscontroller.Config{
		Conn:           natsConn,
		Log:            log,
		MessageService: messageService,
		ConsoleUpdater: consoleUpdater,
	})
	go natsServer.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGKILL)
	<-quit

	if err = natsConn.Close(); err != nil {
		log.Fatal().Err(err).Msg("Closing nats connection error")
	}
}
