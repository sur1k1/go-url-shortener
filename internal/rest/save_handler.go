package rest

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sur1k1/go-url-shortener/internal/util/generate"
	"go.uber.org/zap"
)

type URLSaver interface {
	SaveURL(shortURL string, originalURL string)
}

type SaveHandler struct {
	saver 	URLSaver
	pubAddr string
	log 		*zap.Logger
}

func NewSaveHandler(r *chi.Mux, u URLSaver, pubAddr string, log *zap.Logger) {
	handler := &SaveHandler{
		saver: u,
		pubAddr: pubAddr,
		log: log,
	}

	r.Post("/", handler.SaveHandler)
}

func (h *SaveHandler) SaveHandler(rw http.ResponseWriter, req *http.Request) {
	const op = "rest.SaveHandler"

	// Проверка заголовка на корректность
	// contentType := req.Header.Get("Content-Type")
	// if !strings.HasPrefix(contentType, "text/plain") {
	// 	h.log.Info(
	// 		"incorrect content type",
	// 		zap.String("path", op),
	// 	)

	// 	http.Error(rw, "incorrect content type", http.StatusBadRequest)
	// 	return
	// }
	
	// Чтение тела запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Info(
			"failed to read body",
			zap.String("path", op),
		)

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Валидация запроса
	if len(body) == 0{
		h.log.Info(
			"request body is nil",
			zap.String("path", op),
		)

		http.Error(rw, "body is nil", http.StatusBadRequest)
		return
	}

	// Создание нового URL
	id := generate.GenerateID()

	// Сохранение новой ссылки
	h.saver.SaveURL(id, string(body))
	
	// Формирование ответа клиенту
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	_, err = rw.Write([]byte(h.pubAddr + "/" + id))
	if err != nil {
		h.log.Info(
			"failed to send response",
			zap.String("path", op),
		)

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
