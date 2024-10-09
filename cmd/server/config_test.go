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
				host: "localhost",
				port: 8080,
			},
			"localhost:8080",
			false,
		},
		{
			"127.0.0.1",
			endpoint{
				host: "127.0.0.1",
				port: 8080,
			},
			"127.0.0.1:8080",
			false,
		},
		{
			"127.0.0.1",
			endpoint{
				host: "localhost",
				port: 8080,
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
				host: "localhost",
				port: 8080,
			},
			"localhost:8080",
		},
		{
			"localhost",
			endpoint{
				host: "localhost",
				port: 80,
			},
			"localhost:80",
		},
		{
			"127.0.0.1:80",
			endpoint{
				host: "127.0.0.1",
				port: 80,
			},
			"127.0.0.1:80",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &endpoint{
				host: tt.fields.host,
				port: tt.fields.port,
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
				"ADDRESS": "localhost:8080",
			},
			config{
				endpoint: endpoint{
					host: "localhost",
					port: 8080,
				},
			},
		},
		{
			"127.0.0.1:80",
			map[string]string{
				"ADDRESS": "127.0.0.1:80",
			},
			config{
				endpoint: endpoint{
					host: "127.0.0.1",
					port: 80,
				},
			},
		},
		{
			"empty vars",
			map[string]string{},
			config{
				endpoint: endpoint{
					host: "localhost",
					port: 8080,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("ADDRESS")
			for k, v := range tt.args {
				assert.NoError(t, os.Setenv(k, v))
			}
			assert.EqualValues(t, tt.wantCfg, envConfig(defaultConfig()))
		})
	}
}
