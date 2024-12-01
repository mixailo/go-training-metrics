package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mixailo/go-training-metrics/internal/service/logger"
	"github.com/mixailo/go-training-metrics/internal/service/metrics"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
)

type metricsStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (val float64, ok bool)
	GetCounter(name string) (val int64, ok bool)
	Gauges() map[string]float64
	Counters() map[string]int64

	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type databaseConnection interface {
	Ping() error
}

type storageAware struct {
	Db   databaseConnection
	stor metricsStorage
}

func newStorageAware(metricsStorage metricsStorage) *storageAware {
	return &storageAware{stor: metricsStorage}
}

func (sa *storageAware) updateItemValue(w http.ResponseWriter, r *http.Request) {
	mName := chi.URLParam(r, "name")
	mValue := chi.URLParam(r, "value")
	mType := chi.URLParam(r, "type")

	switch mType {
	case metrics.TypeCounter.String():
		// counter type increments stored value
		convertedValue, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		sa.stor.UpdateCounter(mName, convertedValue)
	case metrics.TypeGauge.String():
		// gauge type replaces stored value
		convertedValue, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		sa.stor.UpdateGauge(mName, convertedValue)
	default:
		// unknown type
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (sa *storageAware) value(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var reqData metrics.Metrics
	var resData metrics.Metrics

	err := dec.Decode(&reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !reqData.IsReadable() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch reqData.MType {
	case metrics.TypeCounter.String():
		stored, ok := sa.stor.GetCounter(reqData.ID)
		if !ok {
			stored = 0
		}
		resData = metrics.Metrics{
			ID:    reqData.ID,
			MType: reqData.MType,
			Delta: &stored,
		}
	case metrics.TypeGauge.String():
		stored, ok := sa.stor.GetGauge(reqData.ID)
		if ok {
			resData = metrics.Metrics{
				ID:    reqData.ID,
				MType: reqData.MType,
				Value: &stored,
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	default:
		// unknown type
		logger.Log.Info("unknown type", zap.String("type", reqData.MType))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.Encode(resData)
}

func (sa *storageAware) update(w http.ResponseWriter, r *http.Request) {
	var data metrics.Metrics

	w.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := dec.Decode(&data)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Warn("error decoding", zap.Error(err))
		return
	}

	if !data.IsWritable() {
		logger.Log.Warn("error data not writable")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch data.MType {
	case metrics.TypeCounter.String():
		// counter type increments stored value
		sa.stor.UpdateCounter(data.ID, *data.Delta)
	case metrics.TypeGauge.String():
		// gauge type updates stored value
		sa.stor.UpdateGauge(data.ID, *data.Value)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.Encode(data)
}

func (sa *storageAware) getItemValue(w http.ResponseWriter, r *http.Request) {
	mName := chi.URLParam(r, "name")
	mType := chi.URLParam(r, "type")

	switch mType {
	case metrics.TypeCounter.String():
		v, ok := sa.stor.GetCounter(mName)
		if ok {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, strconv.FormatInt(v, 10))
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case metrics.TypeGauge.String():
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
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
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

func (sa *storageAware) ping(w http.ResponseWriter, r *http.Request) {
	err := sa.Db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (sa *storageAware) store(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(sa.stor)

	if err != nil {
		return err
	}
	logger.Log.Debug("store", zap.String("path", path))
	return nil
}

func (sa *storageAware) restore(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sa.stor)
	if err != nil {
		return err
	}

	logger.Log.Debug("restore", zap.String("path", path), zap.Int("len gauges", len(sa.stor.Gauges())), zap.Int("len counters", len(sa.stor.Counters())))
	return nil
}
