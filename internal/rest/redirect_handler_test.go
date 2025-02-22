package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
)

func TestHandlers_GetHandler(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		shortURL    string
		httpMethod  string
		wantStatus  int
	}{
		{
			name:        "status code 200",
			originalURL: "https://www.google.com/",
			shortURL:    generateID(),
			httpMethod:  http.MethodGet,
			wantStatus:  http.StatusTemporaryRedirect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewStorage()
			s.URLs[serverURL+tt.shortURL] = tt.originalURL
			h := &RedirectHandler{getter: s}

			req := httptest.NewRequest(tt.httpMethod, fmt.Sprintf("/%s", tt.shortURL), nil)
			rw := httptest.NewRecorder()

			h.GetHandler(rw, req)

			res := rw.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.originalURL, res.Header.Get("Location"))
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
