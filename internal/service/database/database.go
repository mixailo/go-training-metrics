package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Config struct {
	DSN string
}

type Connection interface {
	Ping() error
	Connect(c Config) error
}

func NewConnection(c Config) (Connection, error) {
	var p postgres
	err := p.Connect(c)

	return &p, err
}

type postgres struct {
	db *sql.DB
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

	return nil
}
