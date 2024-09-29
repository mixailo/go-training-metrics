package reporter

import (
	"fmt"
	"github.com/mixailo/go-training-metrics/internal/metrics"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ServerEndpoint struct {
	scheme string
	host   string
	port   int
}

func NewServerEndpoint(scheme string, host string, port int) ServerEndpoint {
	return ServerEndpoint{
		scheme: scheme,
		host:   host,
		port:   port,
	}
}

func (se *ServerEndpoint) String() string {
	return se.scheme + "://" + se.host + ":" + strconv.Itoa(se.port)
}

func (se *ServerEndpoint) CreateURL(path string) string {
	return se.String() + "/" + strings.TrimLeft(path, "/")
}

func SendReport(report metrics.Report, endpoint ServerEndpoint) (err error) {
	fmt.Println(report)
	for _, metric := range report.All() {
		err = sendReportMetric(metric, endpoint)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func sendReportMetric(metric metrics.Data, endpoint ServerEndpoint) (err error) {
	reportURLPath := endpoint.CreateURL(reportPath(metric))
	_, err = http.Post(reportURLPath, "text/plain", nil)

	return
}

func reportPath(metric metrics.Data) (result string) {
	result = "/update/" + metric.CounterType.String() + "/" + url.QueryEscape(metric.Name) + "/" + url.QueryEscape(metric.Value)
	return
}
