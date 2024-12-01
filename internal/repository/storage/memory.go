package storage

import "encoding/json"

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *MemStorage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Gauges   map[string]float64 `json:"Gauges"`
		Counters map[string]int64   `json:"Counters"`
	}{
		Gauges:   m.Gauges(),
		Counters: m.Counters(),
	})
}

func (m *MemStorage) UnmarshalJSON(data []byte) error {
	encoded := struct {
		Gauges   map[string]float64 `json:"Gauges"`
		Counters map[string]int64   `json:"Counters"`
	}{}
	err := json.Unmarshal(data, &encoded)
	if err != nil {
		return err
	}
	m.gauges = encoded.Gauges
	m.counters = encoded.Counters

	return nil
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

func (m *MemStorage) GetGauge(name string) (val float64, ok bool) {
	val, ok = m.gauges[name]
	return
}

func (m *MemStorage) GetCounter(name string) (val int64, ok bool) {
	val, ok = m.counters[name]
	return
}

func NewMemStorage() *MemStorage {
	return &MemStorage{gauges: make(map[string]float64), counters: make(map[string]int64)}
}

func (m *MemStorage) Gauges() map[string]float64 {
	return m.gauges
}

func (m *MemStorage) Counters() map[string]int64 {
	return m.counters
}
