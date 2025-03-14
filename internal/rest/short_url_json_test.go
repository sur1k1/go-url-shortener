package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sur1k1/go-url-shortener/internal/models"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
)

func TestShortenJSONHandler_ShortJSONHandler(t *testing.T) {
	const testPublicAddress = "http://localhost:8080/api/shorten"

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
			NewShortJSONHandler(r, s, testPublicAddress)

			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, contentType, resp := testShortenJSONHandlerRequest(
				t,
				tt.httpMethod,
				testPublicAddress,
				tt.contentType,
				tt.reqBody,
			)

			assert.Equal(t, tt.expectedCode, statusCode, "not equal status code")
			
			assert.NotNil(t, resp, "response body is nil")

			assert.Contains(t, contentType, tt.contentType, "incorrect content type")
		})
	}
}

func testShortenJSONHandlerRequest(t *testing.T, method, path, contentType string, body models.ShortenRequest) (int, string, models.ShortenResponse) {
	request := resty.New().R()

	request.Method = method
	request.URL = path
	request.Body = body
	request.SetHeader("Content-Type", contentType)

	resp, err := request.Send()
	require.NoError(t, err)

	var modelResp models.ShortenResponse
	err = json.Unmarshal(resp.Body(), &modelResp)
	require.NoError(t, err)

	return resp.StatusCode(), resp.Header().Get("Content-Type"), modelResp
}