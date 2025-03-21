package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sur1k1/go-url-shortener/internal/logger"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
)

func TestHandlers_SaveHandler(t *testing.T) {
	log, err := logger.New("info")
	require.NoError(t, err)

	s, err := storage.NewStorage(log, "temp_storage.txt")
	require.NoError(t, err)

	tests := []struct {
		name        string
		contentType string
		httpMethod  string
		originalURL string
		wantStatus  int
	}{
		{
			name:        "status code 201",
			contentType: "text/plain",
			httpMethod:  http.MethodPost,
			originalURL: "https://www.google.com/",
			wantStatus:  201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const publicAddress = "http://localhost:8080/"
			const path = "/"

			r := chi.NewRouter()

			log, err := logger.New("info")
			require.NoError(t, err)

			NewSaveHandler(r, s, publicAddress, log)

			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, contentType, body := testSaveRequest(t, ts, tt.httpMethod, path, tt.contentType, strings.NewReader(tt.originalURL))

			assert.NotNil(t, body, "body is nil")
			
			assert.Equal(t, tt.wantStatus, statusCode, "not equal want and actual status code")

			assert.Contains(t, contentType, "text/plain", "incorrect content type")
		})
	}
}

func testSaveRequest(t *testing.T, ts *httptest.Server, method, path, contentType string, body io.Reader) (int, string, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", contentType)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, resp.Header.Get("Content-Type"), string(respBody)
}