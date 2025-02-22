package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
)

func TestHandlers_postHandler(t *testing.T) {
	tests := []struct {
		name        string
		h           *Handlers
		contentType string
		httpMethod  string
		originalURL string
		wantStatus  int
	}{
		{
			name:        "status code 200",
			h:           &Handlers{storage: storage.NewStorage()},
			contentType: "text/plain",
			httpMethod:  http.MethodPost,
			originalURL: "https://www.google.com/",
			wantStatus:  201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tt.originalURL)))

			req.Header.Set("Content-Type", tt.contentType)

			rw := httptest.NewRecorder()

			tt.h.postHandler(rw, req)

			res := rw.Result()

			body, err := io.ReadAll(res.Body)
			defer res.Body.Close()

			assert.NoError(t, err)

			assert.NotNil(t, body)
			assert.Equal(t, tt.wantStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), "text/plain")
		})
	}
}

func TestHandlers_getHandler(t *testing.T) {
	tests := []struct {
		name string
		h    *Handlers
		originalURL string
		shortURL string
		httpMethod string
		wantStatus int
	}{
		{
			name: "status code 200",
			h: &Handlers{storage: storage.NewStorage()},
			originalURL: "https://www.google.com/",
			shortURL: generateID(),
			httpMethod: http.MethodGet,
			wantStatus: http.StatusTemporaryRedirect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.storage.SaveURL(serverURL+tt.shortURL, tt.originalURL)

			req := httptest.NewRequest(tt.httpMethod, fmt.Sprintf("/%s", tt.shortURL), nil)
			rw := httptest.NewRecorder()

			tt.h.getHandler(rw, req)

			res := rw.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.originalURL, res.Header.Get("Location"))
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
