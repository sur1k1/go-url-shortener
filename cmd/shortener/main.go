package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sur1k1/go-url-shortener/internal/config"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
	"github.com/sur1k1/go-url-shortener/internal/rest"
)

func main() {
	// Getting a configuration
	cf := config.MustGetConfig()
	
	// Storage init
	s := storage.NewStorage()

	if err := http.ListenAndServe(cf.ServerAddress, InitRouter(s, cf.PublicAddress)); err != nil{
		panic(err)
	}
}

func InitRouter(s *storage.MemStorage, pubAddr string) *chi.Mux {
	// Router init
	r := chi.NewRouter()

	// Register handlers
	rest.NewRedirectHandler(r, s)
	rest.NewSaveHandler(r, s, pubAddr)

	return r
}