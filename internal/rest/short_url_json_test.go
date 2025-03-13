package rest

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-resty/resty/v2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"github.com/sur1k1/go-url-shortener/internal/models"
// 	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
// )

// func TestShortenJSONHandler_ShortJSONHandler(t *testing.T) {
// 	const testPublicAddress = "http://localhost:8080/api/shorten"

// 	tests := []struct {
// 		name string
// 		contentType string
// 		httpMethod string
// 		expectedCode int
// 		reqBody models.ShortenRequest
// 	}{
// 		{
// 			name: "status code 201",
// 			contentType: "application/json",
// 			httpMethod: http.MethodPost,
// 			expectedCode: 201,
// 			reqBody: models.ShortenRequest{
// 				URL: "https://google.com/",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := storage.NewStorage()

// 			r := chi.NewRouter()
// 			NewSaveHandler(r, s, testPublicAddress)

// 			ts := httptest.NewServer(r)
// 			defer ts.Close()

// 			testShortenJSONHandlerRequest(
// 				t,
// 				ts,
// 				tt.httpMethod,
// 				testPublicAddress,
// 				tt.contentType,
// 				tt.reqBody,
// 				tt.expectedCode,
// 			)
// 		})
// 	}
// }

// func testShortenJSONHandlerRequest(t *testing.T, ts *httptest.Server, method, path, contentType string, body models.ShortenRequest, expectedCode int) (int, string, string) {
// 	request := resty.New().R()

// 	request.Method = method
// 	request.URL = path
// 	request.Body = body
// 	request.SetHeader("Content-Type", "application/json")

// 	resp, err := request.Send()
// 	require.NoError(t, err)

// 	assert.Equal(t, expectedCode, resp.StatusCode(), "response code didn't match expected")

// 	return resp.StatusCode, resp.Header.Get("Content-Type"), string(respBody)
// }