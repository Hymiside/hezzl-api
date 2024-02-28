package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Hymiside/hezzl-api/pkg/models"
)

type RepositoryPostgres struct {
	db *sql.DB
}

func NewRepositoryPostgres(db *sql.DB) *RepositoryPostgres {
	return &RepositoryPostgres{db: db}
}

func (r *RepositoryPostgres) Goods(ctx context.Context, limit, offset int) (models.GoodsResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.GoodsResponse{}, fmt.Errorf("error to begin transaction: %v", err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(
		ctx, 
		`SELECT 
			id,  
			name,
			description,
			priority,
			removed,
			created_at
		FROM goods LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return models.GoodsResponse{}, fmt.Errorf("error to get goods: %v", err)
	}

	var goods = []models.Good{}
	for rows.Next() {
		var good models.Good
		err = rows.Scan(
			&good.ID,
			&good.Name,
			&good.Description,
			&good.Priority,
			&good.Removed,
			&good.CreatedAt,
		)
		if err != nil {
			return models.GoodsResponse{}, fmt.Errorf("error to scan good: %v", err)
		}
		goods = append(goods, good)
	}

	if err = rows.Err(); err != nil {
		return models.GoodsResponse{}, fmt.Errorf("error rows goods: %v", err)
	}

	totalGoods, err := r.totalGoods(ctx, tx)
	if err != nil {
		return models.GoodsResponse{}, fmt.Errorf("error to get total goods: %v", err)
	}

	totalRemovedGoods, err := r.totalRemovedGoods(ctx, tx)
	if err != nil {
		return models.GoodsResponse{}, fmt.Errorf("error to get total removed goods: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return models.GoodsResponse{}, fmt.Errorf("error to commit transaction: %v", err)
	}

	goodsResponse := models.GoodsResponse{
		Goods: goods,
		Meta: models.Meta{
			Limit:   limit,
			Ofset:   offset,
			Total:   totalGoods,
			Removed: totalRemovedGoods,
		},
	}
	return goodsResponse, nil
}

func (r *RepositoryPostgres) totalGoods(ctx context.Context, tx *sql.Tx) (int, error) {
	row := r.db.QueryRowContext(ctx, `SELECT COUNT(id) FROM goods`)
	if row.Err() != nil {
		return 0, fmt.Errorf("error to get total goods: %v", row.Err())
	}

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("error to scan total goods: %v", err)
	}
	return total, nil
}

func (r *RepositoryPostgres) totalRemovedGoods(ctx context.Context, tx *sql.Tx) (int, error) {
	row := r.db.QueryRowContext(ctx, `SELECT COUNT(id) FROM goods WHERE removed = true`)
	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("error to get total removed goods: %v", err)
	}

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("error to scan total removed goods: %v", err)
	}
	return total, nil
}

func (r *RepositoryPostgres) CreateGood(ctx context.Context, projectID int, name string) (models.Good, error) {
	row := r.db.QueryRowContext(
		ctx, 
		`INSERT INTO 
			goods (project_id, name, priority) 
		VALUES 
			($1, $2, (SELECT COALESCE(MAX(priority), 0) + 1 FROM goods)) 
		RETURNING id`, 
		projectID, name)
	if err := row.Err(); err != nil {
		return models.Good{}, fmt.Errorf("error to create good: %v", err)
	}

	var goodID int
	if err := row.Scan(&goodID); err != nil {
		return models.Good{}, fmt.Errorf("error to scan id: %v", err)
	}

	row = r.db.QueryRowContext(
		ctx, 
		`SELECT 
			id, 
			project_id,
			name, 
			description, 
			priority, 
			removed, 
			created_at 
		FROM goods 
		WHERE id = $1`, 
		goodID)
	if err := row.Err(); err != nil {
		return models.Good{}, fmt.Errorf("error to get good: %v", err)
	}

	var good models.Good
	if err := row.Scan(
		&good.ID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreatedAt,
	); err != nil {
		return models.Good{}, fmt.Errorf("error to scan good: %v", err)
	}

	return good, nil
}

func (r *RepositoryPostgres) UpdateGood(ctx context.Context, good models.Good, goodID, projectID int) (models.Good, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Good{}, err
	}

	row := tx.QueryRowContext(
		ctx,
		`UPDATE 
			goods 
		SET 
			name = CASE WHEN $1 = '' THEN name ELSE $1 END, 
			description = CASE WHEN $2 = '' THEN description ELSE $2 END,
			priority = CASE WHEN $3 = 0 THEN priority ELSE $3 END,
			removed = CASE WHEN $4 = false THEN removed ELSE $4 END
		WHERE id = $5 AND project_id = $6 RETURNING id`,
		good.Name,
		good.Description,
		good.Priority,
		good.Removed,
		goodID,
		projectID)
	if err := row.Err(); err != nil {
		return models.Good{}, fmt.Errorf("error to update good: %v", err)
	}

	var id int
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Good{}, fmt.Errorf("good not found: %v", err)
		}
		return models.Good{}, fmt.Errorf("error to scan id: %v", err)
	}

	row = tx.QueryRowContext(
		ctx,
		`SELECT 
			id, 
			project_id,
			name, 
			description, 
			priority, 
			removed, 
			created_at 
		FROM goods 
		WHERE id = $1`, 
		id)
	if err := row.Err(); err != nil {
		return models.Good{}, fmt.Errorf("error to get good: %v", err)
	}

	var goodUpdated models.Good
	if err := row.Scan(
		&goodUpdated.ID,
		&goodUpdated.ProjectID,
		&goodUpdated.Name,
		&goodUpdated.Description,
		&goodUpdated.Priority,
		&goodUpdated.Removed,
		&goodUpdated.CreatedAt,
	); err != nil {
		return models.Good{}, fmt.Errorf("error to scan good: %v", err)
	}
	
	if err = tx.Commit(); err != nil {
		return models.Good{}, fmt.Errorf("error to commit transaction: %v", err)
	}
	return goodUpdated, nil
}