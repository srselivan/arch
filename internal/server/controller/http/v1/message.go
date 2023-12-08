package v1

import (
	"arch/internal/server/entity"
	"arch/pkg/uuid"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func (h *Handler) InitMessageRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.createMessage)

	return r
}

func (h *Handler) createMessage(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("io.ReadAll: %v", err), http.StatusBadRequest)
		return
	}
	h.log.Debug().Msgf("Request payload: %s", payload)

	message := entity.Message{
		ID:          uuid.NewV7(),
		MediaObject: nil,
		TextObject: &entity.TextObject{
			Body: string(payload),
		},
	}
	h.messageService.Save(r.Context(), message)

	w.WriteHeader(http.StatusCreated)
}
