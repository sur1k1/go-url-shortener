package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
)

func TestHandlers_SaveHandler(t *testing.T) {
	tests := []struct {
		name        string
		h           *SaveHandler
		contentType string
		httpMethod  string
		originalURL string
		wantStatus  int
	}{
		{
			name:        "status code 200",
			h:           &SaveHandler{saver: storage.NewStorage()},
			contentType: "text/plain",
			httpMethod:  http.MethodPost,
			originalURL: "https://www.google.com/",
			wantStatus:  201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.originalURL))

			req.Header.Set("Content-Type", tt.contentType)

			rw := httptest.NewRecorder()

			tt.h.SaveHandler(rw, req)

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