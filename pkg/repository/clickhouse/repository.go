package clickhouse

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hymiside/hezzl-api/pkg/models"
)

type RepositoryClickhouse struct {
	db *sql.DB
}

func NewRepositoryClickhouse(db *sql.DB) *RepositoryClickhouse {
	return &RepositoryClickhouse{db: db}
}

func (r *RepositoryClickhouse) CreateLogs(ctx context.Context, logs []models.Good) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error to begin transaction: %v", err)
	}
	defer tx.Rollback()

	for _, v := range logs {

		fmt.Printf("v: %v\n", v)
		if _, err := tx.ExecContext(
			ctx, 
			`INSERT INTO 
				logs (id, project_id, name, description, priority, removed, created_at) 
			VALUES
				(?, ?, ?, ?, ?, ?, now())`, 
			v.ID, 
			v.ProjectID, 
			v.Name, 
			v.Description, 
			v.Priority, 
			v.Removed, 
		); err != nil {
			return fmt.Errorf("error to create logs: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error to commit transaction: %v", err)
	}
	return nil
}