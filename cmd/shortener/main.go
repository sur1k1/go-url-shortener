package main

import (
	"net/http"

	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
	"github.com/sur1k1/go-url-shortener/internal/rest"
)

func main() {
	// Storage init
	s := storage.NewStorage()

	// Server init
	mux := http.NewServeMux()

	// Register handlers
	rest.NewRedirectHandler(mux, s)
	rest.NewSaveHandler(mux, s)

	if err := http.ListenAndServe(`:8080`, mux); err != nil{
		panic(err)
	}
}