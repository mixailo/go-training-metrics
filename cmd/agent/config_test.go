package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEndpoint_Parse(t *testing.T) {
	tests := []struct {
		name            string
		desiredEndpoint Endpoint
		input           string
		wantErr         bool
	}{
		{
			"localhost",
			Endpoint{
				Host: "localhost",
				Port: 8080,
			},
			"localhost:8080",
			false,
		},
		{
			"127.0.0.1",
			Endpoint{
				Host: "127.0.0.1",
				Port: 8080,
			},
			"127.0.0.1:8080",
			false,
		},
		{
			"127.0.0.1",
			Endpoint{
				Host: "localhost",
				Port: 8080,
			},
			"127.0.0.18080",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := Endpoint{}
			if tt.wantErr {
				assert.Error(t, tt.desiredEndpoint.Parse(tt.input))
			} else {
				assert.NoError(t, ep.Parse(tt.input))
				assert.EqualValues(t, tt.desiredEndpoint, ep)
			}
		})
	}
}

func TestEndpoint_String(t *testing.T) {
	tests := []struct {
		name   string
		fields Endpoint
		want   string
	}{
		{
			"localhost:8080",
			Endpoint{
				Host: "localhost",
				Port: 8080,
			},
			"localhost:8080",
		},
		{
			"localhost",
			Endpoint{
				Host: "localhost",
				Port: 80,
			},
			"localhost:80",
		},
		{
			"127.0.0.1:80",
			Endpoint{
				Host: "127.0.0.1",
				Port: 80,
			},
			"127.0.0.1:80",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Endpoint{
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
		wantCfg Config
	}{
		{
			"default",
			map[string]string{
				"ADDRESS":         "localhost:8080",
				"REPORT_INTERVAL": "10",
				"POLL_INTERVAL":   "2",
			},
			Config{
				Endpoint: Endpoint{
					Host: "localhost",
					Port: 8080,
				},
				ReportInterval: 10,
				PollInterval:   2,
			},
		},
		{
			"127.0.0.1:80",
			map[string]string{
				"ADDRESS":         "127.0.0.1:80",
				"REPORT_INTERVAL": "10",
				"POLL_INTERVAL":   "2",
			},
			Config{
				Endpoint: Endpoint{
					Host: "127.0.0.1",
					Port: 80,
				},
				ReportInterval: 10,
				PollInterval:   2,
			},
		},
		{
			"127.0.0.1:80 with intervals",
			map[string]string{
				"ADDRESS":         "127.0.0.1:80",
				"REPORT_INTERVAL": "100",
				"POLL_INTERVAL":   "20",
			},
			Config{
				Endpoint: Endpoint{
					Host: "127.0.0.1",
					Port: 80,
				},
				ReportInterval: 100,
				PollInterval:   20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.args {
				assert.NoError(t, os.Setenv(k, v))
			}
			assert.EqualValues(t, tt.wantCfg, EnvConfig(DefaultConfig()))
		})
	}
}
