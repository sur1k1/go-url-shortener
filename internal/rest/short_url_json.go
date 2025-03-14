package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sur1k1/go-url-shortener/internal/models"
	"github.com/sur1k1/go-url-shortener/internal/util/generate"
	"go.uber.org/zap"
)

type ShortenJSONHandler struct {
	saver 	URLSaver // интерфейс описан в save_handler.go
	pubAddr string
	log 		*zap.Logger
}

func NewShortJSONHandler(r *chi.Mux, u URLSaver, pubAddr string, log *zap.Logger) {
	handler := &ShortenJSONHandler{
		saver:   u,
		pubAddr: pubAddr,
		log: log,
	}

	r.Post("/api/shorten", handler.ShortJSONHandler)
}

func (h *ShortenJSONHandler) ShortJSONHandler(rw http.ResponseWriter, req *http.Request) {
	const op = "rest.ShortJSONHandler"

	// Проверка заголовка на корректность
	contentType := req.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		h.log.Info(
			"incorrect content type",
			zap.String("path", op),
		)

		http.Error(rw, "incorrect content type", http.StatusBadRequest)
		return
	}

	// Чтение тела запроса
	var reqBody models.ShortenRequest

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		h.log.Info(
			"failed to read body",
			zap.String("path", op),
		)

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


	respBody, err := json.Marshal(resp)
	if err != nil {
		h.log.Info(
			"failed to marshal body",
			zap.String("path", op),
		)

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(respBody)
	if err != nil {
		h.log.Info(
			"failed to send response",
			zap.String("path", op),
		)

		return
	}
}