package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

const (
	logsPath = "_logs"
)

func New(level, serviceName string) *zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	output := zerolog.ConsoleWriter{
		TimeFormat: time.RFC3339Nano,
		Out:        os.Stdout,
	}

	if _, err = os.Stat(logsPath); os.IsNotExist(err) {
		if err = os.Mkdir(logsPath, 0777); err != nil {
			panic(err)
		}
	}

	logsFilePath := fmt.Sprintf("%s/%s.log", logsPath, serviceName)
	file, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}

	multi := zerolog.MultiLevelWriter(output, file)

	l := zerolog.New(multi).With().Caller().Timestamp().Logger()
	l.Level(logLevel)

	return &l
}
