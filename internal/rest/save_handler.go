package rest

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"
)

const serverURL = "http://localhost:8080/"

type URLSaver interface {
	SaveURL(shortURL string, originalURL string)
}

type SaveHandler struct {
	saver URLSaver
}

func NewSaveHandler(mux *http.ServeMux, u URLSaver) {
	handler := &SaveHandler{
		saver: u,
	}

	mux.HandleFunc("/", handler.PostHandler)
}

func (h *SaveHandler) PostHandler(rw http.ResponseWriter, req *http.Request) {
	// Проверка метода запроса
	if req.Method != http.MethodPost{
		http.Error(rw, "incorrect request method", http.StatusBadRequest)
		return
	}

	// Проверка заголовка на корректность
	contentType := req.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/plain") {
		http.Error(rw, "incorrect content type", http.StatusBadRequest)
		return
	}

	// Чтение тела запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация запроса
	if len(body) == 0{
		http.Error(rw, "body is nil", http.StatusBadRequest)
		return
	}

	// Создание нового URL
	replacedURL := serverURL + generateID()

	// Сохранение новой ссылки
	h.saver.SaveURL(replacedURL, string(body))

	// Формирование ответа клиенту
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write([]byte(replacedURL))
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