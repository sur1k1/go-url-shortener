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

	// Server init
	r := chi.NewRouter()

	// Register handlers
	rest.NewRedirectHandler(r, s)
	rest.NewSaveHandler(r, s)

	if err := http.ListenAndServe(`:8080`, r); err != nil{
		panic(err)
	}
}