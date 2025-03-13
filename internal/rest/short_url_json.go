package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sur1k1/go-url-shortener/internal/models"
	"github.com/sur1k1/go-url-shortener/internal/util/generate"
)

type ShortenJSONHandler struct {
	saver 	URLSaver // интерфейс описан в save_handler.go
	pubAddr string
}

func NewShortJSONHandler(r *chi.Mux, u URLSaver, pubAddr string) {
	handler := &ShortenJSONHandler{
		saver:   u,
		pubAddr: pubAddr,
	}

	r.Post("/api/shorten", handler.ShortJSONHandler)
}

func (h *ShortenJSONHandler) ShortJSONHandler(rw http.ResponseWriter, req *http.Request) {
	// Проверка заголовка на корректность
	contentType := req.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(rw, "incorrect content type", http.StatusBadRequest)
		return
	}

	// Чтение тела запроса
	var reqBody models.ShortenRequest

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(rw, "failed to get body", http.StatusBadRequest)
		return
	}

	// Создание нового URL
	id := generate.GenerateID()

	// Сохранение новой ссылки
	h.saver.SaveURL(id, reqBody.URL)

	resp := models.ShortenResponse{
		Reslut: h.pubAddr + "/" + id,
	}

	// Формирование ответа клиенту
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)


	if err := json.NewEncoder(rw).Encode(&resp); err != nil {
		http.Error(rw, "error encoding response", http.StatusInternalServerError)
		return
	}
}