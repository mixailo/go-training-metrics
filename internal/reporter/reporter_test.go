package reporter

import (
	"github.com/mixailo/go-training-metrics/internal/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServerEndpoint(t *testing.T) {
	endpoint := NewServerEndpoint("http", "localhost", 8080)
	assert.IsType(t, endpoint, ServerEndpoint{})
}

func TestSendReport(t *testing.T) {
	t.Skip()
}

func TestServerEndpoint_CreateUrl(t *testing.T) {
	endpoint := NewServerEndpoint("http", "localhost", 8080)

	assert.Equal(t, "http://localhost:8080/path", endpoint.CreateUrl("/path"))
	assert.Equal(t, "http://localhost:8080/path", endpoint.CreateUrl("path"))

}

func TestServerEndpoint_String(t *testing.T) {
	var endpoint ServerEndpoint
	endpoint = NewServerEndpoint("http", "localhost", 8080)
	assert.Equal(t, "http://localhost:8080", endpoint.String())
	endpoint = NewServerEndpoint("https", "localhost", 8080)
	assert.Equal(t, "https://localhost:8080", endpoint.String())
}

func Test_reportPath(t *testing.T) {
	type args struct {
		metric metrics.Data
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "counter",
			args: args{
				metric: metrics.Data{
					Name:        "test",
					Value:       "13444",
					CounterType: metrics.TypeCounter,
				},
			},
			wantResult: "/update/counter/test/13444",
		},
		{
			name: "gauge",
			args: args{
				metric: metrics.Data{
					Name:        "test",
					Value:       "13444",
					CounterType: metrics.TypeGauge,
				},
			},
			wantResult: "/update/gauge/test/13444",
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.wantResult, reportPath(tt.args.metric))
	}
}

func Test_sendReportMetric(t *testing.T) {
	t.Skip()
}
