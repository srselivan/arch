package v1

import (
	"arch/internal/server/entity"
	"arch/pkg/uuid"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	fileDir = "_upload/"
)

func (h *Handler) InitFileRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.createFile)

	return r
}

func (h *Handler) createFile(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("r.FormFile: %v", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Error().Err(err).Msg("An error occurred while closing file")
		}
	}()

	ext := filepath.Ext(fileHeader.Filename)
	switch ext {
	case ".jpg", ".png", ".jpeg":
	default:
		http.Error(w, "File extension is not allowed", http.StatusBadRequest)
		return
	}

	filePath := fileDir + uuid.NewV7() + ext
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("os.Create: %v", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err = dst.Close(); err != nil {
			log.Error().Err(err).Msg("An error occurred while closing file")
		}
	}()
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, fmt.Sprintf("io.Copy: %v", err), http.StatusInternalServerError)
		return
	}

	message := entity.Message{
		ID: uuid.NewV7(),
		MediaObject: &entity.MediaObject{
			Filename: fileHeader.Filename,
			Path:     filePath,
		},
		TextObject: nil,
	}
	h.messageService.Save(r.Context(), message)

	w.WriteHeader(http.StatusCreated)
}
