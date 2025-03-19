package rest

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sur1k1/go-url-shortener/internal/models"
	"github.com/sur1k1/go-url-shortener/internal/repository"
	"go.uber.org/zap"
)

type URLGetter interface {
	GetURL(shortURL string) (*models.URLData, error)
}

type RedirectHandler struct {
	getter URLGetter
	log 	*zap.Logger
}

func NewRedirectHandler(r *chi.Mux, u URLGetter, log *zap.Logger) {
	handler := &RedirectHandler{
		getter: u,
		log: log,
	}

	r.Get("/{id}", handler.RedirectHandler)
}

func (h *RedirectHandler) RedirectHandler(rw http.ResponseWriter, req *http.Request) {
	const op = "rest.RedirectHandler"

	// Парсинг URL для получения ID
	id := req.URL.Path[1:]
	
	// Поиск ID в базе данных
	urlData, err := h.getter.GetURL(id)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			h.log.Info(
				"url not found",
				zap.String("path", op),
				zap.String("id", id),
			)

			http.Error(rw, "url not found", http.StatusNotFound)
			return
		}

		h.log.Info(
			"failed to get url",
			zap.String("path", op),
			zap.String("id", id),
		)

		http.Error(rw, "failed to get url", http.StatusInternalServerError)
		return
	}

	// Формирование ответа клиенту
	rw.Header().Set("Location", urlData.OriginalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}