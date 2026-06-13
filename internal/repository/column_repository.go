package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/model"
)

type ColumnRepository struct {
	db *pgxpool.Pool
}

func NewColumnRepository(db *pgxpool.Pool) *ColumnRepository {
	return &ColumnRepository{db: db}
}

func (r *ColumnRepository) Create(ctx context.Context, col *model.Column) error {
	query := `
		INSERT INTO columns (id, board_id, name, position, color)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		col.ID, col.BoardID, col.Name, col.Position, col.Color,
	).Scan(&col.CreatedAt, &col.UpdatedAt)
}

func (r *ColumnRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Column, error) {
	query := `
		SELECT id, board_id, name, position, color, created_at, updated_at
		FROM columns WHERE id = $1`

	col := &model.Column{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&col.ID, &col.BoardID, &col.Name, &col.Position, &col.Color,
		&col.CreatedAt, &col.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get column: %w", err)
	}
	return col, nil
}

func (r *ColumnRepository) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]model.Column, error) {
	query := `
		SELECT id, board_id, name, position, color, created_at, updated_at
		FROM columns WHERE board_id = $1
		ORDER BY position ASC`

	rows, err := r.db.Query(ctx, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}
	defer rows.Close()

	var columns []model.Column
	for rows.Next() {
		var col model.Column
		if err := rows.Scan(&col.ID, &col.BoardID, &col.Name, &col.Position, &col.Color, &col.CreatedAt, &col.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan column: %w", err)
		}
		columns = append(columns, col)
	}
	return columns, nil
}

func (r *ColumnRepository) Update(ctx context.Context, col *model.Column) error {
	query := `
		UPDATE columns SET name = $2, position = $3, color = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query, col.ID, col.Name, col.Position, col.Color).Scan(&col.UpdatedAt)
}

func (r *ColumnRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM columns WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *ColumnRepository) Reorder(ctx context.Context, columnIDs []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	for i, id := range columnIDs {
		_, err := tx.Exec(ctx, `UPDATE columns SET position = $2, updated_at = NOW() WHERE id = $1`, id, i)
		if err != nil {
			return fmt.Errorf("reorder column %d: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *ColumnRepository) GetMaxPosition(ctx context.Context, boardID uuid.UUID) (int, error) {
	var maxPos int
	query := `SELECT COALESCE(MAX(position), -1) FROM columns WHERE board_id = $1`
	err := r.db.QueryRow(ctx, query, boardID).Scan(&maxPos)
	return maxPos, err
}
