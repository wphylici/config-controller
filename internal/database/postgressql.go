package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var Psql *PostgreSQL

type PostgreSQL struct {
	config                  *Config
	db                      *sql.DB
	serviceConfigRepository *ServiceConfigRepository
}

func New(config *Config) *PostgreSQL {
	return &PostgreSQL{
		config: config,
	}
}

func (p *PostgreSQL) Open() error {
	db, err := sql.Open("postgres", p.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	p.db = db

	return nil
}

func (p *PostgreSQL) Close() {
	p.db.Close()
}

func (p *PostgreSQL) ServiceConfig() *ServiceConfigRepository {
	if p.serviceConfigRepository != nil {
		return p.serviceConfigRepository
	}

	p.serviceConfigRepository = &ServiceConfigRepository{
		psql: p,
	}

	return p.serviceConfigRepository
}
