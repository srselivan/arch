package monitor

import (
	"arch/internal/monitor/service"
	"arch/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

const (
	serviceName = "monitor"
	watchDir    = "_heartbeats"
)

func Run() {
	log := logger.New("debug", serviceName)

	runner := service.NewRunner(log)
	if err := runner.RunInstance(); err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msg("Runner successfully start new server instance")

	watchdog := service.NewWatchdog(watchDir, log, runner)
	log.Info().Msgf("Run watchdog for %s directory", watchDir)

	if err := watchdog.Run(); err != nil {
		log.Fatal().Err(err).Send()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGKILL)
	<-quit

	log.Info().Msg("Shutdown watchdog process")
	watchdog.Down()

	log.Info().Msg("Stop monitor")
}
