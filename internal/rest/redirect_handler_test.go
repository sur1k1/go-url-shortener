package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
	"github.com/sur1k1/go-url-shortener/internal/util/generate"
)

func TestHandlers_RedirectHandler(t *testing.T) {
	const tempURL = "http://localhost:8080/"

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
			s := storage.NewStorage()
			s.URLs[tt.shortURL] = tt.originalURL

			h := &RedirectHandler{getter: s}

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