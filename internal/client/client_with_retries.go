package client

import (
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

const RetryUntilSuccess = -1

type ClientWithRetries struct {
	client http.Client
	logger *zerolog.Logger
}

func NewClientWithRetries(logger *zerolog.Logger) *ClientWithRetries {
	return &ClientWithRetries{
		logger: logger,
	}
}

func (r *ClientWithRetries) Do(request *http.Request, retryTimeout time.Duration, retriesCount int) {
	for {
		if retriesCount == 0 {
			return
		}
		retriesCount--

		response, err := r.client.Do(request)
		if err == nil {
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				r.logger.Error().Err(err).Msg("Response body read error")
				return
			}

			r.logger.Debug().
				Str("status", response.Status).
				Str("body", string(body)).
				Send()
			return
		}

		r.logger.Error().
			Err(err).
			Func(func(e *zerolog.Event) {
				if retriesCount > RetryUntilSuccess {
					e.Int("retry left", retriesCount)
				}
			}).
			Msg("Do http request error")

		select {
		case <-request.Context().Done():
			return
		case <-time.After(retryTimeout):
		}
	}
}
