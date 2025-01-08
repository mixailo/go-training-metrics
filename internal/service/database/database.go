package database

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/mixailo/go-training-metrics/internal/service/logger"
)

type Config struct {
	DSN string
}

const counterTable = "counter_storage"
const gaugeTable = "gauge_storage"

type Connection interface {
	Ping() error
	Connect(c Config) error
	Close() error
	HasTables() bool
	CreateTables() error
	Gauges() map[string]float64
	Counters() map[string]int64
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (val float64, ok bool)
	GetCounter(name string) (val int64, ok bool)
}

func NewConnection(c Config) (Connection, error) {
	var p postgres
	err := p.Connect(c)

	return &p, err
}

type postgres struct {
	db  *sql.DB
	dsn string
}

//go:embed migrations/*.sql
var fs embed.FS

func (p *postgres) HasTables() bool {
	if !p.hasTable(counterTable) {
		return false
	}
	if !p.hasTable(gaugeTable) {
		return false
	}
	return true
}

func (p *postgres) hasTable(table string) bool {
	var result bool
	row := p.db.QueryRow(`select exists(select from information_schema.tables where table_schema='public' and table_name=$1);`, table)
	err := row.Scan(&result)
	if err != nil {
		return false
	}

	return result
}

func (p *postgres) Counters() map[string]int64 {
	result := map[string]int64{}

	rows, err := p.db.Query(`SELECT * from counter_storage`)

	if err != nil {
		logger.Log.Error(err.Error())
		return result
	}

	for rows.Next() {
		var name string
		var value int64
		if err := rows.Scan(&name, &value); err == nil {
			result[name] = value
		}
	}

	return result
}

func (p *postgres) Gauges() map[string]float64 {
	result := map[string]float64{}

	rows, err := p.db.Query(`SELECT * from gauge_storage`)

	if err != nil {
		logger.Log.Error(err.Error())
		return result
	}

	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err == nil {
			result[name] = value
		}
	}

	return result
}

func (p *postgres) CreateTables() error {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, p.dsn)
	if err != nil {
		return err
	}
	err = m.Up()

	// no change means no errors
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}

func (p *postgres) Ping() error {
	return p.db.Ping()
}

func (p *postgres) Connect(c Config) error {
	db, err := sql.Open("postgres", c.DSN)
	if err != nil {
		return err
	}

	p.db = db
	p.dsn = c.DSN

	return nil
}

func (p *postgres) Close() error {
	return p.db.Close()
}

func (p *postgres) GetGauge(name string) (val float64, ok bool) {
	row := p.db.QueryRow(`SELECT * from gauge_storage where value_name=$1;`, name)
	err := row.Scan(&val)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false
	}

	return val, true
}

func (p *postgres) GetCounter(name string) (val int64, ok bool) {
	row := p.db.QueryRow(`SELECT * from counter_storage where value_name=$1;`, name)
	err := row.Scan(&val)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false
	}

	return val, true
}

func (p *postgres) UpdateGauge(name string, value float64) {
	_, err := p.db.Exec(`INSERT INTO gauge_storage as t (value_name, gauge_value) VALUES ($1, $2) ON CONFLICT (value_name) DO UPDATE SET gauge_value=$3;`, name, value, value)
	if err != nil {
		logger.Log.Error("error updating gauge",
			zap.Error(err),
			zap.String("name", name),
			zap.Float64("value", value),
		)
	}
}

func (p *postgres) UpdateCounter(name string, value int64) {
	_, err := p.db.Exec(`INSERT INTO counter_storage as t (value_name, counter_value) VALUES ($1, $2) ON CONFLICT (value_name) DO UPDATE SET counter_value=t.counter_value+$3;`, name, value, value)
	if err != nil {
		logger.Log.Error("error updating counter",
			zap.Error(err),
			zap.String("name", name),
			zap.Int64("value", value),
		)
	}
}
