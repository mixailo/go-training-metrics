package sender

import (
	"bytes"
	"compress/gzip"
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

	// create gzip encoder
	zl, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return err
	}

	// json-encode request body
	enc := json.NewEncoder(zl)
	err = enc.Encode(metric)

	if err != nil {
		return err
	}
	zl.Close()

	// send request
	reportURLPath := endpoint.CreateURL(reportPath())
	request, err := http.NewRequest("POST", reportURLPath, &buf)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")

	response, err := http.DefaultClient.Do(request)

	if err == nil {
		defer response.Body.Close()
	}

	return err
}

func reportPath() (result string) {
	result = "/update/"
	return
}
