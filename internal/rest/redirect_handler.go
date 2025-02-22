package rest

import "net/http"

type URLGetter interface {
	GetURL(shortURL string) (string, bool)
}

type RedirectHandler struct {
	getter URLGetter
}

func NewRedirectHandler(mux *http.ServeMux, u URLGetter) {
	handler := &RedirectHandler{
		getter: u,
	}

	mux.HandleFunc("/{id}", handler.GetHandler)
}

func (h *RedirectHandler) GetHandler(rw http.ResponseWriter, req *http.Request) {
	// Проверка метода запроса
	if req.Method != http.MethodGet {
		http.Error(rw, "incorrect request method", http.StatusBadRequest)
		return
	}

	// Валидация запроса
	if len(req.URL.Path) < 1 {
		http.Error(rw, "incorrect url", http.StatusBadRequest)
		return
	}

	// Парсинг URL для получения ID
	id := req.URL.Path[1:]

	// Поиск ID в базе данных
	originalURL, ok := h.getter.GetURL(serverURL + id)
	if !ok {
		http.Error(rw, "id not found", http.StatusBadRequest)
		return
	}

	// Формирование ответа клиенту
	rw.Header().Set("Location", originalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}