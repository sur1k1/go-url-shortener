package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"

	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
)

const serverURL = "http://localhost:8080/"

// Описание интерфейса базы данных
type Storage interface {
	GetURL(shortURL string) (string, bool)
	SaveURL(shortURL string, originalURL string)
}

// Структура для хэндлеров с полем доступа к сторэджу
type Handlers struct {
	storage Storage
}

func main() {
	// Storage init
	s := storage.NewStorage()

	// Server init
	mux := http.NewServeMux()

	// Register handlers
	h := Handlers{storage: s}
	mux.HandleFunc("/", h.postHandler)
	mux.HandleFunc("/{id}", h.getHandler)

	if err := http.ListenAndServe(`:8080`, mux); err != nil{
		panic(err)
	}
}

func (h *Handlers) postHandler(rw http.ResponseWriter, req *http.Request) {
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
	h.storage.SaveURL(replacedURL, string(body))

	// Формирование ответа клиенту
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write([]byte(replacedURL))
	if err != nil {
		log.Println("cannot send response", err)
	}
}

func (h *Handlers) getHandler(rw http.ResponseWriter, req *http.Request) {
	// Проверка метода запроса
	if req.Method != http.MethodGet{
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
	originalURL, ok := h.storage.GetURL(serverURL+id)
	if !ok {
		http.Error(rw, "id not found", http.StatusBadRequest)
		return
	}

	// Формирование ответа клиенту
	rw.Header().Set("Location", originalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}