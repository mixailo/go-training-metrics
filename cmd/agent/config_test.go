package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpoint_Parse(t *testing.T) {
	tests := []struct {
		name            string
		desiredEndpoint endpoint
		input           string
		wantErr         bool
	}{
		{
			"localhost",
			endpoint{
				Host: "localhost",
				Port: 8080,
			},
			"localhost:8080",
			false,
		},
		{
			"127.0.0.1",
			endpoint{
				Host: "127.0.0.1",
				Port: 8080,
			},
			"127.0.0.1:8080",
			false,
		},
		{
			"127.0.0.1",
			endpoint{
				Host: "localhost",
				Port: 8080,
			},
			"127.0.0.18080",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoint{}
			if tt.wantErr {
				assert.Error(t, tt.desiredEndpoint.parse(tt.input))
			} else {
				assert.NoError(t, ep.parse(tt.input))
				assert.EqualValues(t, tt.desiredEndpoint, ep)
			}
		})
	}
}

func TestEndpoint_String(t *testing.T) {
	tests := []struct {
		name   string
		fields endpoint
		want   string
	}{
		{
			"localhost:8080",
			endpoint{
				Host: "localhost",
				Port: 8080,
			},
			"localhost:8080",
		},
		{
			"localhost",
			endpoint{
				Host: "localhost",
				Port: 80,
			},
			"localhost:80",
		},
		{
			"127.0.0.1:80",
			endpoint{
				Host: "127.0.0.1",
				Port: 80,
			},
			"127.0.0.1:80",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &endpoint{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			assert.Equal(t, tt.want, e.String())
		})
	}
}

func TestEnvConfig(t *testing.T) {
	type args map[string]string

	tests := []struct {
		name    string
		args    args
		wantCfg config
	}{
		{
			"default",
			map[string]string{
				"ADDRESS":         "localhost:8080",
				"REPORT_INTERVAL": "10",
				"POLL_INTERVAL":   "2",
				"LOG_LEVEL":       "info",
			},
			config{
				endpoint: endpoint{
					Host: "localhost",
					Port: 8080,
				},
				reportInterval: 10,
				pollInterval:   2,
				logLevel:       "info",
			},
		},
		{
			"127.0.0.1:80",
			map[string]string{
				"ADDRESS":         "127.0.0.1:80",
				"REPORT_INTERVAL": "10",
				"POLL_INTERVAL":   "2",
				"LOG_LEVEL":       "info",
			},
			config{
				endpoint: endpoint{
					Host: "127.0.0.1",
					Port: 80,
				},
				reportInterval: 10,
				pollInterval:   2,
				logLevel:       "info",
			},
		},
		{
			"127.0.0.1:80 with intervals",
			map[string]string{
				"ADDRESS":         "127.0.0.1:80",
				"REPORT_INTERVAL": "100",
				"POLL_INTERVAL":   "20",
				"LOG_LEVEL":       "info",
			},
			config{
				endpoint: endpoint{
					Host: "127.0.0.1",
					Port: 80,
				},
				reportInterval: 100,
				pollInterval:   20,
				logLevel:       "info",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("ADDRESS")
			os.Unsetenv("REPORT_INTERVAL")
			os.Unsetenv("POLL_INTERVAL")
			os.Unsetenv("LOG_LEVEL")
			for k, v := range tt.args {
				assert.NoError(t, os.Setenv(k, v))
			}
			assert.EqualValues(t, tt.wantCfg, envConfig(defaultConfig()))
		})
	}
}
