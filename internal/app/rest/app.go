package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sur1k1/go-url-shortener/internal/config"
	"github.com/sur1k1/go-url-shortener/internal/models"
	"github.com/sur1k1/go-url-shortener/internal/rest"
	"github.com/sur1k1/go-url-shortener/internal/rest/middlewares"
	"go.uber.org/zap"
)

type ServiceRepository interface {
	GetURL(shortURL string) (models.URLData, bool)
	SaveURL(urlData models.URLData)
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
	cp := middlewares.NewCompressMiddleware(log)

	// Register middlewares
	r.Use(lm.Logger, cp.Compress)
	
	// Register handlers
	rest.NewRedirectHandler(r, repo, log)
	rest.NewSaveHandler(r, repo, cf.BaseURL, log)
	rest.NewShortJSONHandler(r, repo, cf.BaseURL, log)

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