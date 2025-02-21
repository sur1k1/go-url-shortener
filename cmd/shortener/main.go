package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"sync"
)

var memstorage = map[string]string{}
const serverURL = "http://localhost:8080/"
var c sync.RWMutex

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", postHandler)
	mux.HandleFunc("/{id}", getHandler)

	if err := http.ListenAndServe(`:8080`, mux); err != nil{
		panic(err)
	}
}

func postHandler(rw http.ResponseWriter, req *http.Request) {
	// Проверка метода запроса
	if req.Method != http.MethodPost{
		http.Error(rw, "incorrect request method", http.StatusBadRequest)
		return
	}

	// Проверка заголовка на корректность
	if req.Header.Get("Content-Type") != "text/plain"{
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
	c.Lock()
	memstorage[replacedURL] = string(body)
	c.Unlock()

	// Формирование ответа клиенту
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(replacedURL))
}

func getHandler(rw http.ResponseWriter, req *http.Request) {
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

	c.RLock()
	// Поиск ID в базе данных
	originalURL, ok := memstorage[serverURL+id]
	if !ok {
		http.Error(rw, "id not found", http.StatusBadRequest)
		return
	}
	c.RUnlock()

	// Формирование ответа клиенту
	rw.Header().Set("Location", originalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}