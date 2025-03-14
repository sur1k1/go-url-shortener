package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type URLGetter interface {
	GetURL(shortURL string) (string, bool)
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

	// Валидация запроса
	if len(req.URL.Path) < 1 {
		h.log.Info(
			"invalid url path",
			zap.String("path", op),
		)

		http.Error(rw, "incorrect url", http.StatusBadRequest)
		return
	}

	// Парсинг URL для получения ID
	id := req.URL.Path[1:]
	
	// Поиск ID в базе данных
	originalURL, ok := h.getter.GetURL(id)
	if !ok {
		h.log.Info(
			"id not found",
			zap.String("path", op),
			zap.String("url", originalURL),
		)

		http.Error(rw, "id not found", http.StatusNotFound)
		return
	}

	// Формирование ответа клиенту
	rw.Header().Set("Location", originalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}