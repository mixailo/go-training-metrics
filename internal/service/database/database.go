package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
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

	if err != nil {
		return nil, err
	}

	return &p, nil
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
