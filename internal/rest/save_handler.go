package rest

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type URLSaver interface {
	SaveURL(shortURL string, originalURL string)
}

type SaveHandler struct {
	saver 	URLSaver
	pubAddr string
}

func NewSaveHandler(r *chi.Mux, u URLSaver, pubAddr string) {
	handler := &SaveHandler{
		saver: u,
		pubAddr: pubAddr,
	}

	r.Post("/", handler.SaveHandler)
}

func (h *SaveHandler) SaveHandler(rw http.ResponseWriter, req *http.Request) {
	// Проверка заголовка на корректность
	contentType := req.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/plain") {
		http.Error(rw, "incorrect content type", http.StatusBadRequest)
		return
	}

	// Чтение тела запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Валидация запроса
	if len(body) == 0{
		http.Error(rw, "body is nil", http.StatusBadRequest)
		return
	}

	// Создание нового URL
	id := generateID()

	// Сохранение новой ссылки
	h.saver.SaveURL(id, string(body))
	
	// Формирование ответа клиенту
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write([]byte(h.pubAddr + "/" + id))
	if err != nil {
		log.Println("cannot send response", err)
		return
	}
}

func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}