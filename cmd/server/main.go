package main

import (
	"github.com/go-chi/chi"
	"github.com/mixailo/go-training-metrics/internal/metrics"
	"github.com/mixailo/go-training-metrics/internal/storage"
	"io"
	"log"
	"net/http"
	"strconv"
)

type MetricsStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (val float64, ok bool)
	GetCounter(name string) (val int64, ok bool)
	Gauges() map[string]float64
	Counters() map[string]int64
}

type storageAware struct {
	stor MetricsStorage
}

func newStorageAware() *storageAware {
	s := storage.NewStorage()
	return &storageAware{stor: &s}
}

func (sa *storageAware) updateItemValue(w http.ResponseWriter, r *http.Request) {
	mName := chi.URLParam(r, "name")
	mValue := chi.URLParam(r, "value")
	mType := chi.URLParam(r, "type")

	if mType == metrics.TypeCounter.String() {
		// counter type increments stored value
		convertedValue, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sa.stor.UpdateCounter(mName, convertedValue)
	} else if mType == metrics.TypeGauge.String() {
		// gauge type replaces stored value
		convertedValue, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sa.stor.UpdateGauge(mName, convertedValue)
	} else {
		// unknown type
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (sa *storageAware) getItemValue(w http.ResponseWriter, r *http.Request) {
	mName := chi.URLParam(r, "name")
	mType := chi.URLParam(r, "type")

	if mType == metrics.TypeCounter.String() {
		v, ok := sa.stor.GetCounter(mName)
		if ok {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, strconv.FormatInt(v, 10))
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else if mType == metrics.TypeGauge.String() {
		v, ok := sa.stor.GetGauge(mName)
		if ok {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, strconv.FormatFloat(v, 'f', -1, 64))
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (sa *storageAware) getAllValues(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	head := `<html><head><title>All Metrics</title></head><body><table>`
	io.WriteString(w, head)
	for k, v := range sa.stor.Gauges() {
		io.WriteString(w, "<tr><td>Gauge</td><td>"+k+"</td><td>"+strconv.FormatFloat(v, 'f', -1, 64)+"</td></tr>")
	}
	for k, v := range sa.stor.Counters() {
		io.WriteString(w, "<tr><td>Counter</td><td>"+k+"</td><td>"+strconv.FormatInt(v, 10)+"</td></tr>")
	}
	foot := `</table></body></html>`
	io.WriteString(w, foot)
}

func main() {
	ParseFlags()
	sa := newStorageAware()
	router := chi.NewRouter()

	router.Post("/update/{type}/{name}/{value}", sa.updateItemValue)
	router.Get("/value/{type}/{name}", sa.getItemValue)
	router.Get("/", sa.getAllValues)

	log.Fatal(http.ListenAndServe(serverEndpoint.String(), router))

}
