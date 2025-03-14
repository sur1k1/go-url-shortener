package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sur1k1/go-url-shortener/internal/logger"
	"github.com/sur1k1/go-url-shortener/internal/models"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
)

func TestShortenJSONHandler_ShortJSONHandler(t *testing.T) {
	const publicAddress = "http://localhost:8080/"
	const path = "/api/shorten"

	tests := []struct {
		name string
		contentType string
		httpMethod string
		expectedCode int
		reqBody models.ShortenRequest
	}{
		{
			name: "status code 201",
			contentType: "application/json",
			httpMethod: http.MethodPost,
			expectedCode: 201,
			reqBody: models.ShortenRequest{
				URL: "https://google.com/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewStorage()

			r := chi.NewRouter()

			log, err := logger.New("info")
			require.NoError(t, err)

			NewShortJSONHandler(r, s, publicAddress, log)

			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, contentType, resp := testShortenJSONHandlerRequest(
				t,
				ts,
				tt.httpMethod,
				path,
				tt.contentType,
				tt.reqBody,
			)

			assert.Equal(t, tt.expectedCode, statusCode, "not equal status code")
			
			assert.NotNil(t, resp, "response body is nil")

			assert.Contains(t, contentType, tt.contentType, "incorrect content type")
		})
	}
}

func testShortenJSONHandlerRequest(t *testing.T, ts *httptest.Server, method, path, contentType string, body models.ShortenRequest) (int, string, models.ShortenResponse) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(body)
	require.NoError(t, err, "failed to encode body")

	req, err := http.NewRequest(method, ts.URL+path, &buf)
	require.NoError(t, err, "failed to create request")

	req.Header.Set("Content-Type", contentType)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err, "failed to send request")
	defer resp.Body.Close()

	var modelResp models.ShortenResponse
	err = json.NewDecoder(resp.Body).Decode(&modelResp)
	require.NoError(t, err, "failed to encode resp body")

	return resp.StatusCode, resp.Header.Get("Content-Type"), modelResp
}