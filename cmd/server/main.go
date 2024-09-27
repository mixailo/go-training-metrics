package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.gauges[name] = value
}

func (m *MemStorage) UpdateCounter(name string, value int64) {
	oldValue, ok := m.counters[name]
	if !ok {
		m.counters[name] = value
	} else {
		m.counters[name] = oldValue + value
	}
}

func validatePost(r *http.Request) (status int, err error) {
	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, errors.New("invalid method, only POST is allowed")
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		return http.StatusUnsupportedMediaType, errors.New("only text/plain is supported")
	}

	return
}

func update(w http.ResponseWriter, r *http.Request) {
	status, err := validatePost(r)
	if err != nil {
		w.WriteHeader(status)
		return
	}

	fragments := strings.Split(r.URL.Path, "/")

	if len(fragments) != 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	counterType := fragments[2]
	name := fragments[3]
	unConvertedValue := fragments[4]

	if len(name) < 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if counterType == TypeCounter {
		value, err := strconv.ParseInt(unConvertedValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		memStorage.UpdateCounter(name, value)
	} else if counterType == TypeGauge {
		value, err := strconv.ParseFloat(unConvertedValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		memStorage.UpdateGauge(name, value)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

const (
	TypeCounter string = "counter"
	TypeGauge   string = "gauge"
)

func NewMemStorage() MemStorage {
	return MemStorage{gauges: make(map[string]float64), counters: make(map[string]int64)}
}

var memStorage = NewMemStorage()

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", update)
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
