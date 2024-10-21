package sender

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mixailo/go-training-metrics/internal/service/logger"
	"github.com/mixailo/go-training-metrics/internal/service/metrics"
)

type ServerEndpoint struct {
	Scheme string
	Host   string
	Port   int
}

func NewServerEndpoint(scheme string, host string, port int) ServerEndpoint {
	return ServerEndpoint{
		Scheme: scheme,
		Host:   host,
		Port:   port,
	}
}

func (se *ServerEndpoint) String() string {
	return se.Scheme + "://" + se.Host + ":" + strconv.Itoa(se.Port)
}

func (se *ServerEndpoint) CreateURL(path string) string {
	return se.String() + "/" + strings.TrimLeft(path, "/")
}

func SendReport(report metrics.Report, endpoint ServerEndpoint) (err error) {
	for _, metric := range report.All() {
		logger.Log.Info("send report", zap.String("metric", metric.String()))
		err = sendReportMetricWithRetries(metric, endpoint)
		if err != nil {
			logger.Log.Info("error", zap.Error(err))
		}
	}
	return
}

func sendReportMetricWithRetries(metric metrics.Metrics, endpoint ServerEndpoint) (err error) {
	for i := 0; i < 3; i++ {
		err = sendReportMetric(metric, endpoint)
		if err == nil {
			break
		} else {
			logger.Log.Info("error, will retry request after 0.05 secs", zap.Error(err))
			time.Sleep(50 * time.Millisecond)
		}
	}

	return err
}

func sendReportMetric(metric metrics.Metrics, endpoint ServerEndpoint) error {
	var buf bytes.Buffer

	// encode request body
	enc := json.NewEncoder(&buf)
	err := enc.Encode(metric)
	if err != nil {
		return err
	}

	// send request
	reportURLPath := endpoint.CreateURL(reportPath())
	response, err := http.Post(reportURLPath, "application/json", &buf)
	if err == nil {
		defer response.Body.Close()
	}

	return err
}

func reportPath() (result string) {
	result = "/update/"
	return
}
