package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mixailo/go-training-metrics/internal/repository/storage"
	"github.com/mixailo/go-training-metrics/internal/service/logger"
	"github.com/mixailo/go-training-metrics/internal/service/metrics"
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

type storageAware struct {
	stor metricsStorage
}

func newStorageAware(metricsStorage metricsStorage) *storageAware {
	return &storageAware{stor: metricsStorage}
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
	reqDataJSON, _ := json.Marshal(reqData)
	logger.Log.Info("data", zap.String("request", string(reqDataJSON)))
	logger.Log.Info("data", zap.String("storage", fmt.Sprint(sa.stor)))
	if reqData.MType == metrics.TypeCounter.String() {
		// counter type increments stored value
		updated, ok := sa.stor.GetCounter(reqData.ID)
		if !ok {
			updated = 0
		}
		resData = metrics.Metrics{
			ID:    reqData.ID,
			MType: reqData.MType,
			Delta: &updated,
		}
	} else if reqData.MType == metrics.TypeGauge.String() {
		updated, ok := sa.stor.GetGauge(reqData.ID)
		if ok {
			resData = metrics.Metrics{
				ID:    reqData.ID,
				MType: reqData.MType,
				Value: &updated,
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}

	} else {
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

	reqDataJSON, _ := json.Marshal(data)
	logger.Log.Info("data", zap.String("request", string(reqDataJSON)))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Info("error decoding", zap.Error(err))
		return
	}

	if !data.IsWritable() {
		logger.Log.Info("error data not writable")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data.MType == metrics.TypeCounter.String() {
		// counter type increments stored value
		sa.stor.UpdateCounter(data.ID, *data.Delta)
	} else if data.MType == metrics.TypeGauge.String() {
		sa.stor.UpdateGauge(data.ID, *data.Value)
	} else {
		// unknown type
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

func newMux(sa *storageAware) *chi.Mux {
	router := chi.NewRouter()

	router.Use(gzipMiddleware)
	router.Use(logger.RequestResponseLogger)

	router.Post("/update/{type}/{name}/{value}", sa.updateItemValue)
	router.Get("/value/{type}/{name}", sa.getItemValue)
	router.Post("/update/", sa.update)
	router.Post("/value/", sa.value)
	router.Get("/", sa.getAllValues)

	return router
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
	logger.Log.Info("store", zap.String("path", path))
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

	logger.Log.Info("restore", zap.String("path", path), zap.Int("len gauges", len(sa.stor.Gauges())), zap.Int("len counters", len(sa.stor.Counters())))
	return nil
}

// yet ungraceful
func gracefulShutdownCatcher(path string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Log.Info("shutdown signal caught")
		shutdown(path)
	}()
}

func persistenceTicker(interval int64, path string) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	logger.Log.Debug("init persistence ticker", zap.Int64("interval", interval))

	for {
		select {
		case <-ticker.C:
			err := sa.store(path)
			if err != nil {
				logger.Log.Error("persistence ticker error", zap.Error(err), zap.String("path", path))
			} else {
				logger.Log.Debug("data stored by ticking timer")
			}
		}
	}
}

func shutdown(path string) {
	logger.Log.Info("shutting down gracefully, let's save data to disk", zap.String("path", path))
	sa.store(path)
	os.Exit(0)
}

var sa *storageAware

func main() {
	// init logging
	serverConf, err := initConfig()
	if err != nil {
		panic(err) // logger has not been initialized yet
	}
	if err := logger.Initialize(serverConf.logLevel); err != nil {
		panic(err) // cannot log without logger
	}

	gracefulShutdownCatcher(serverConf.fileStoragePath)

	// init storage
	sa = newStorageAware(storage.NewMemStorage())
	if serverConf.doRestoreValues {
		sa.restore(serverConf.fileStoragePath)
	}

	logger.Log.Info(fmt.Sprintf("Starting server at %s:%d", serverConf.endpoint.host, serverConf.endpoint.port))

	chiMux := newMux(sa)
	if serverConf.storeInterval == 0 {
		logger.Log.Info("will save data to disk immediately")
		chiMux.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
				err := sa.store(serverConf.fileStoragePath)
				if err != nil {
					logger.Log.Error(err.Error())
				}
			})
		})
	} else {
		logger.Log.Info("will save data to disk periodically", zap.Int64("interval", serverConf.storeInterval))
		go persistenceTicker(serverConf.storeInterval, serverConf.fileStoragePath)
		logger.Log.Info("init ok")
	}

	err = http.ListenAndServe(serverConf.endpoint.String(), chiMux)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
