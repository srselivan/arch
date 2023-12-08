package client

import (
	"arch/pkg/logger"
	"context"
	"encoding/base64"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	serviceName         = "client"
	defaultRetryTimeout = 10 * time.Second
)

type credentials struct {
	login, password string
}

func (c *credentials) base64() string {
	s := c.login + ":" + c.password
	return base64.StdEncoding.EncodeToString([]byte(s))
}

var (
	retriesCount int
	retryTimeout time.Duration
	payload      string
	creds        credentials
)

func Run() {
	flag.DurationVar(&retryTimeout, "timeout", defaultRetryTimeout, "request sending retry timeout")
	flag.IntVar(&retriesCount, "retries", RetryUntilSuccess, "request sending retry count")
	flag.StringVar(&payload, "m", "", "request payload")
	flag.StringVar(&creds.login, "l", "root", "login")
	flag.StringVar(&creds.password, "p", "root", "password")
	flag.Parse()

	log := logger.New("debug", serviceName)

	client := NewClientWithRetries(log)

	requestContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		request *http.Request
		err     error
	)

	if isFilePath(payload) {
		request, err = GetRequestWithFile(requestContext, payload)
		if err != nil {
			log.Error().Err(err).Msg("Creating request error")
		}
	} else {
		request, err = GetRequestWithTextBody(requestContext, payload)
		if err != nil {
			log.Error().Err(err).Msg("Creating request error")
		}
	}

	out := make(chan struct{})
	go func() {
		client.Do(request, retryTimeout, retriesCount)
		close(out)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGKILL)

	select {
	case <-quit:
		cancel()
	case <-out:
	}
}

func isFilePath(filePath string) bool {
	if "" == filepath.Ext(filePath) {
		return false
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
