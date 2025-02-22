package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLGetter interface {
	GetURL(shortURL string) (string, bool)
}

type RedirectHandler struct {
	getter URLGetter
}

func NewRedirectHandler(r *chi.Mux, u URLGetter) {
	handler := &RedirectHandler{
		getter: u,
	}

	r.Get("/{id}", handler.RedirectHandler)
}

func (h *RedirectHandler) RedirectHandler(rw http.ResponseWriter, req *http.Request) {
	// Валидация запроса
	if len(req.URL.Path) < 1 {
		http.Error(rw, "incorrect url", http.StatusBadRequest)
		return
	}

	// Парсинг URL для получения ID
	id := req.URL.Path[1:]
	
	// Поиск ID в базе данных
	originalURL, ok := h.getter.GetURL(id)
	if !ok {
		http.Error(rw, "id not found", http.StatusBadRequest)
		return
	}

	// Формирование ответа клиенту
	rw.Header().Set("Location", originalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}