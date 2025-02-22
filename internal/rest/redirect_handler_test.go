package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
)

func TestHandlers_RedirectHandler(t *testing.T) {
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

			r := chi.NewRouter()
			NewRedirectHandler(r, s)

			ts := httptest.NewServer(r)
			defer ts.Close()

			resp := testRedirectRequest(t, ts, tt.httpMethod, fmt.Sprintf("/%s", tt.shortURL))

			assert.Equal(t, tt.originalURL, resp.Header.Get("Location"))
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func testRedirectRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}