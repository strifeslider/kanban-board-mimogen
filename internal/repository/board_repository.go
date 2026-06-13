package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/model"
)

type BoardRepository struct {
	db *pgxpool.Pool
}

func NewBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(ctx context.Context, board *model.Board) error {
	query := `
		INSERT INTO boards (id, workspace_id, name, description, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		board.ID, board.WorkspaceID, board.Name, board.Description, board.CreatedBy,
	).Scan(&board.CreatedAt, &board.UpdatedAt)
}

func (r *BoardRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Board, error) {
	query := `
		SELECT id, workspace_id, name, description, created_by, created_at, updated_at
		FROM boards WHERE id = $1 AND deleted_at IS NULL`

	board := &model.Board{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&board.ID, &board.WorkspaceID, &board.Name, &board.Description,
		&board.CreatedBy, &board.CreatedAt, &board.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}
	return board, nil
}

func (r *BoardRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.Board, error) {
	query := `
		SELECT id, workspace_id, name, description, created_by, created_at, updated_at
		FROM boards WHERE workspace_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list boards: %w", err)
	}
	defer rows.Close()

	var boards []model.Board
	for rows.Next() {
		var b model.Board
		if err := rows.Scan(&b.ID, &b.WorkspaceID, &b.Name, &b.Description, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan board: %w", err)
		}
		boards = append(boards, b)
	}
	return boards, nil
}

func (r *BoardRepository) Update(ctx context.Context, board *model.Board) error {
	query := `
		UPDATE boards SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query, board.ID, board.Name, board.Description).Scan(&board.UpdatedAt)
}

func (r *BoardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE boards SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *BoardRepository) AddMember(ctx context.Context, member *model.BoardMember) error {
	query := `
		INSERT INTO board_members (id, board_id, user_id, role)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`

	return r.db.QueryRow(ctx, query,
		member.ID, member.BoardID, member.UserID, member.Role,
	).Scan(&member.CreatedAt)
}

func (r *BoardRepository) RemoveMember(ctx context.Context, boardID, userID uuid.UUID) error {
	query := `DELETE FROM board_members WHERE board_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, boardID, userID)
	return err
}

func (r *BoardRepository) IsMember(ctx context.Context, boardID, userID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM board_members WHERE board_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, query, boardID, userID).Scan(&exists)
	return exists, err
}
