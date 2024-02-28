package postgres

import "database/sql"

type RepositoryPostgres struct {
	db *sql.DB
}

func NewRepositoryPostgres(db *sql.DB) *RepositoryPostgres {
	return &RepositoryPostgres{db: db}
}