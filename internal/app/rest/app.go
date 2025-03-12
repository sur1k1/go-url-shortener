package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sur1k1/go-url-shortener/internal/config"
	"github.com/sur1k1/go-url-shortener/internal/rest"
	"github.com/sur1k1/go-url-shortener/internal/rest/middlewares"
	"go.uber.org/zap"
)

type ServiceRepository interface {
	GetURL(shortURL string) (string, bool)
	SaveURL(shortURL string, originalURL string)
}

type App struct {
	log *zap.Logger
	r 	*chi.Mux
	cf 	*config.Config
}

func New(log *zap.Logger, repo ServiceRepository, cf *config.Config) *App {
	// Router init
	r := chi.NewRouter()

	// Init middlewares
	lm := middlewares.NewLoggerMiddleware(log)

	// Register middlewares
	r.Use(lm.Logger)
	
	// Register handlers
	rest.NewRedirectHandler(r, repo)
	rest.NewSaveHandler(r, repo, cf.BaseURL)

	return &App{
		log: log,
		r: r,
		cf: cf,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Panic("run server", zap.Error(err))
	}
}

func (a *App) Run() error {
	const op = "rest.Run"

	if err := http.ListenAndServe(a.cf.ServerAddress, a.r); err != nil{
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}