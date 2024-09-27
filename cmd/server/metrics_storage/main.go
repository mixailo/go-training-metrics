package metrics_storage

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

func NewStorage() MetricsStorage {
	return &MemStorage{gauges: make(map[string]float64), counters: make(map[string]int64)}
}

type MetricsStorage interface {
	UpdateCounter(name string, value int64)
	UpdateGauge(name string, value float64)
}
