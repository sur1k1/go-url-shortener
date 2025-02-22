package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
	"github.com/sur1k1/go-url-shortener/internal/rest"
)

func main() {
	// Storage init
	s := storage.NewStorage()

	if err := http.ListenAndServe(`:8080`, InitRouter(s)); err != nil{
		panic(err)
	}
}

func InitRouter(s *storage.MemStorage) *chi.Mux {
	// Router init
	r := chi.NewRouter()

	// Register handlers
	rest.NewRedirectHandler(r, s)
	rest.NewSaveHandler(r, s)

	return r
}