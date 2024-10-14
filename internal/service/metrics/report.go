package metrics

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Type string

const (
	TypeCounter Type = "counter"
	TypeGauge   Type = "gauge"
)

func (t Type) String() string {
	return string(t)
}

type Report struct {
	value map[string]Metrics
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// check if Metrics is ok for storing
func (m *Metrics) IsWritable() bool {
	if m.ID == "" {
		return false
	}
	if m.MType == TypeCounter.String() && m.Delta != nil {
		return true
	} else if m.MType == TypeGauge.String() && m.Value != nil {
		return true
	}

	return false
}

// check if Metrics is ok for reading
func (m *Metrics) IsReadable() bool {
	if m.ID == "" {
		return false
	}
	if m.MType == TypeCounter.String() || m.MType == TypeGauge.String() {
		return true
	}
	return false
}

func (m *Metrics) String() string {
	decoded, err := json.Marshal(m)
	if err != nil {
		return ""
	} else {
		return string(decoded)
	}
}

func NewReport() Report {
	r := Report{}
	r.value = make(map[string]Metrics)
	return r
}

func (r *Report) All() []Metrics {
	result := make([]Metrics, 0, len(r.value))
	for _, v := range r.value {
		result = append(result, v)
	}

	return result
}

func (r *Report) Length() int {
	return len(r.value)
}

func (r *Report) Get(name string) (result Metrics, ok bool) {
	result, ok = r.value[name]
	return
}

func (r *Report) Has(name string) bool {
	_, ok := r.value[name]
	return ok
}

func NewMetrics(counterType Type, name, value string) (res Metrics, err error) {
	res.ID = name
	if counterType == TypeCounter {
		var v int64
		res.MType = TypeCounter.String()
		v, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return
		}
		res.Delta = &v
	} else if counterType == TypeGauge {
		var v float64
		res.MType = TypeGauge.String()
		v, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return
		}
		res.Value = &v
	} else {
		err = errors.New("invalid metric type")
	}

	return
}

func (r *Report) Add(metric Metrics) {
	r.value[metric.ID] = metric
}

func (r *Report) AddUnConverted(counterType Type, name, value string) error {
	m, err := NewMetrics(counterType, name, value)
	if err == nil {
		r.Add(m)
	}
	return err
}
