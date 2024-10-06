package storage

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

func (m *MemStorage) GetGauge(name string) (val float64, ok bool) {
	val, ok = m.gauges[name]
	return
}

func (m *MemStorage) GetCounter(name string) (val int64, ok bool) {
	val, ok = m.counters[name]
	return
}

func NewStorage() *MemStorage {
	return &MemStorage{gauges: make(map[string]float64), counters: make(map[string]int64)}
}

func (m *MemStorage) Gauges() map[string]float64 {
	return m.gauges
}

func (m *MemStorage) Counters() map[string]int64 {
	return m.counters
}
