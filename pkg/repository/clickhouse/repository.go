package clickhouse

import "database/sql"

type RepositoryClickhouse struct {
	db *sql.DB
}

func NewRepositoryClickhouse(db *sql.DB) *RepositoryClickhouse {
	return &RepositoryClickhouse{db: db}
}