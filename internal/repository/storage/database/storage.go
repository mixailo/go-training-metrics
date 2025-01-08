package database

import "github.com/mixailo/go-training-metrics/internal/service/database"

type dbStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (val float64, ok bool)
	GetCounter(name string) (val int64, ok bool)
	Gauges() map[string]float64
	Counters() map[string]int64
	Ping() error
}

type Storage struct {
	db database.Connection
}

func NewStorage(db database.Connection) (*Storage, error) {
	var storage *Storage

	err := db.Ping()
	if err == nil {
		storage = &Storage{
			db: db,
		}
		err = storage.CreateTables()
	}

	return storage, err
}

func (s *Storage) CreateTables() error {
	var err error

	if !s.db.HasTables() {
		err = s.db.CreateTables()
	}

	return err
}

func (s *Storage) Gauges() map[string]float64 {
	return s.db.Gauges()
}

func (s *Storage) Counters() map[string]int64 {
	return s.db.Counters()
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}

func (s *Storage) UpdateGauge(name string, value float64) {
	s.db.UpdateGauge(name, value)
}

func (s *Storage) UpdateCounter(name string, value int64) {
	s.db.UpdateCounter(name, value)
}

func (s *Storage) GetGauge(name string) (val float64, ok bool) {
	return s.db.GetGauge(name)
}

func (s *Storage) GetCounter(name string) (val int64, ok bool) {
	return s.db.GetCounter(name)
}
