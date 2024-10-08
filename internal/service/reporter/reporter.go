package reporter

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
		err = sendReportMetric(metric, endpoint)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func sendReportMetric(metric metrics.Data, endpoint ServerEndpoint) error {
	reportURLPath := endpoint.CreateURL(reportPath(metric))
	response, err := http.Post(reportURLPath, "text/plain", nil)
	if err == nil {
		defer response.Body.Close()
	}

	return err
}

func reportPath(metric metrics.Data) (result string) {
	result = "/update/" + metric.CounterType.String() + "/" + url.QueryEscape(metric.Name) + "/" + url.QueryEscape(metric.Value)
	return
}
