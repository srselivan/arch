package service

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

const (
	beatInterval = 10 * time.Second
)

type Watchdog interface {
	Run() error
	Down()
}

type watchdog struct {
	logger         *zerolog.Logger
	heartbeatsPath string
	down           chan struct{}
	runner         Runner
}

func NewWatchdog(heartbeatsPath string, logger *zerolog.Logger, runner Runner) Watchdog {
	return &watchdog{
		logger:         logger,
		heartbeatsPath: heartbeatsPath,
		down:           make(chan struct{}),
		runner:         runner,
	}
}

func (w *watchdog) Run() error {
	stat, err := os.Stat(w.heartbeatsPath)
	if err != nil {
		return fmt.Errorf("watchdog run: os.Stat: %w", err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("heartbeats path (%s) is not a directory", w.heartbeatsPath)
	}

	initialStates := make(map[string]time.Time, 1)
	ticker := time.NewTicker(beatInterval)
	checker := func() {
		files, err := os.ReadDir(w.heartbeatsPath)
		if err != nil {
			w.logger.Err(err).Send()
		}

		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				w.logger.Err(err).Send()
			}

			modified, ok := initialStates[info.Name()]
			if !ok {
				w.logger.Info().Msgf("File (%s) not found as trackable, adding to it", info.Name())
				initialStates[info.Name()] = info.ModTime()
				continue
			}

			if !modified.Before(info.ModTime()) {
				service, _ := strings.CutSuffix(info.Name(), ".txt")
				w.logger.Warn().Msgf("File size not changed, dont have heartbeat from service (%s)", service)
				w.logger.Info().Msgf("Trying restart service (%s)", service)

				delete(initialStates, info.Name())
				if err := os.Remove(w.heartbeatsPath + "/" + info.Name()); err != nil {
					w.logger.Error().Err(err).Msgf("Remove file (%s) has been failed", info.Name())
				}

				if err := w.runner.Reboot(service); err != nil {
					w.logger.Error().Err(err).Msgf("Reboot service (%s) has been failed", service)
					continue
				}

				w.logger.Info().Err(err).Msgf("Successfully reboot service (%s)", service)
				continue
			}

			initialStates[info.Name()] = info.ModTime()
		}
	}

	go func() {
		for {
			select {
			case <-w.down:
				return
			case <-ticker.C:
				checker()
			}
		}
	}()

	return nil
}

func (w *watchdog) Down() {
	w.down <- struct{}{}
}
