package sender

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mixailo/go-training-metrics/internal/service/metrics"
)

func TestNewServerEndpoint(t *testing.T) {
	endpoint := NewServerEndpoint("http", "localhost", 8080)
	assert.IsType(t, ServerEndpoint{}, endpoint)
}

func TestSendReport(t *testing.T) {
	t.Skip()
}

func TestServerEndpoint_CreateUrl(t *testing.T) {
	endpoint := NewServerEndpoint("http", "localhost", 8080)

	assert.Equal(t, "http://localhost:8080/path", endpoint.CreateURL("/path"))
	assert.Equal(t, "http://localhost:8080/path", endpoint.CreateURL("path"))

}

func TestServerEndpoint_String(t *testing.T) {
	var endpoint ServerEndpoint
	endpoint = NewServerEndpoint("http", "localhost", 8080)
	assert.Equal(t, "http://localhost:8080", endpoint.String())
	endpoint = NewServerEndpoint("https", "localhost", 8080)
	assert.Equal(t, "https://localhost:8080", endpoint.String())
}

func Test_reportPath(t *testing.T) {
	cv := int64(10)
	type args struct {
		metric metrics.Metrics
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "counter",
			args: args{
				metric: metrics.Metrics{
					ID:    "test",
					Delta: &cv,
					MType: metrics.TypeCounter.String(),
				},
			},
			wantResult: "/update/",
		},
		{
			name: "gauge",
			args: args{
				metric: metrics.Metrics{
					ID:    "test",
					Delta: &cv,
					MType: metrics.TypeCounter.String(),
				},
			},
			wantResult: "/update/",
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.wantResult, reportPath())
	}
}

func Test_sendReportMetric(t *testing.T) {
	t.Skip()
}
