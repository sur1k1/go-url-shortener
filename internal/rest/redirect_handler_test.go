package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sur1k1/go-url-shortener/internal/config"
	"github.com/sur1k1/go-url-shortener/internal/logger"
	"github.com/sur1k1/go-url-shortener/internal/models"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
	"github.com/sur1k1/go-url-shortener/internal/util/generate"
)

func TestHandlers_RedirectHandler(t *testing.T) {
	const tempURL = "http://localhost:8080/"

	// Getting a configuration
	cf := config.MustGetConfig()

	log, err := logger.New("info")
	require.NoError(t, err)

	s, err := storage.NewStorage(log, cf.FilePath)
	require.NoError(t, err)

	tests := []struct {
		name        string
		originalURL string
		shortURL    string
		httpMethod  string
		wantStatus  int
	}{
		{
			name:        "status code 307",
			originalURL: "https://stackoverflow.com/",
			shortURL:    generate.GenerateID(),
			httpMethod:  http.MethodGet,
			wantStatus:  http.StatusTemporaryRedirect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.SaveURL(&models.URLData{
				UUID: "1",
				ShortURL: tt.shortURL,
				OriginalURL: tt.originalURL,
			})
			require.NoError(t, err)
			
			log, err := logger.New("info")
			require.NoError(t, err)

			h := &RedirectHandler{getter: s, log: log}

			req := httptest.NewRequest(tt.httpMethod, tempURL+tt.shortURL, nil)
			rw := httptest.NewRecorder()

			h.RedirectHandler(rw, req)

			resp := rw.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.originalURL, resp.Header.Get("Location"), "not equal original and actual url in header")
			assert.Equal(t, tt.wantStatus, resp.StatusCode, "not equal want and actual status code")
		})
	}
}