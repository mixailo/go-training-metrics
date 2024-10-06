package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/mixailo/go-training-metrics/internal/repository/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_newStorageAware(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"valid type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.IsType(t, newStorageAware(storage.NewStorage()), &storageAware{})
		})
	}
}

func Test_storageAware_getAllValues(t *testing.T) {
	sa := newStorageAware(storage.NewStorage())

	server := httptest.NewServer(http.HandlerFunc(sa.getAllValues))

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
}

func Test_storageAware_getItemValue(t *testing.T) {

	type fields struct {
		stor MetricsStorage
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &storageAware{
				stor: tt.fields.stor,
			}
			sa.getItemValue(tt.args.w, tt.args.r)
		})
	}
}

func Test_storageAware_updateItemValue(t *testing.T) {
	type fields struct {
		stor MetricsStorage
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &storageAware{
				stor: tt.fields.stor,
			}
			sa.updateItemValue(tt.args.w, tt.args.r)
		})
	}
}
