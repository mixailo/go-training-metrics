package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mixailo/go-training-metrics/internal/repository/storage"
)

func Test_newStorageAware(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"valid type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.IsType(t, newStorageAware(storage.NewMemStorage()), &storageAware{})
		})
	}
}

func Test_storageAware_getAllValues(t *testing.T) {
	sa := newStorageAware(storage.NewMemStorage())

	server := httptest.NewServer(newMux(sa))

	defer server.Close()

	tests := []struct {
		name         string
		url          string
		status       int
		hasSubstring string
	}{
		{
			"title page",
			"/",
			http.StatusOK,
			"All Metrics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = http.MethodGet
			req.URL = server.URL + tt.url

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.status, resp.StatusCode(), "unexpected status code")
			if len(tt.hasSubstring) > 0 {
				assert.True(t, strings.Contains(string(resp.Body()), tt.hasSubstring), "expected substring not found")
			}
		})
	}

	t.Run("gzip compression", func(t *testing.T) {
		req := resty.New().SetDoNotParseResponse(true).R()
		req.Method = http.MethodGet
		req.Header.Add("Accept-Encoding", "gzip")
		req.URL = server.URL + "/"

		resp, err := req.Send()
		assert.NoError(t, err, "error making HTTP request")
		zr, err := gzip.NewReader(resp.RawBody())
		require.NoError(t, err, "error creating gzip reader")

		assert.Equal(t, http.StatusOK, resp.StatusCode(), "unexpected status code")
		require.True(t, strings.Contains(resp.Header().Get("Content-Encoding"), "gzip"), "no Content-Encoding header or no gzip in it")

		_, err = io.ReadAll(zr)
		require.NoError(t, err)
		defer resp.RawResponse.Body.Close()
	})
}
