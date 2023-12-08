package v1

import (
	"arch/internal/server/controller/http/template"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) InitWebRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.webPage)

	return r
}

func (h *Handler) webPage(w http.ResponseWriter, r *http.Request) {
	messages, err := h.messageService.GetAll(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("h.messageService.GetAll: %v", err), http.StatusInternalServerError)
		return
	}

	templateContent := template.TemplateContentFromServiceModels(messages)
	h.log.Debug().Any("content", templateContent).Send()

	if err = h.template.Execute(w, templateContent); err != nil {
		http.Error(w, fmt.Sprintf("h.template.Execute: %v", err), http.StatusInternalServerError)
		return
	}
}
