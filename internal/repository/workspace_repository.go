package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/model"
)

type WorkspaceRepository struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepository(db *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) Create(ctx context.Context, ws *model.Workspace) error {
	query := `
		INSERT INTO workspaces (id, name, slug, description, owner_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		ws.ID, ws.Name, ws.Slug, ws.Description, ws.OwnerID,
	).Scan(&ws.CreatedAt, &ws.UpdatedAt)
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error) {
	query := `
		SELECT id, name, slug, description, owner_id, created_at, updated_at
		FROM workspaces WHERE id = $1 AND deleted_at IS NULL`

	ws := &model.Workspace{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ws.ID, &ws.Name, &ws.Slug, &ws.Description, &ws.OwnerID,
		&ws.CreatedAt, &ws.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get workspace: %w", err)
	}
	return ws, nil
}

func (r *WorkspaceRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Workspace, error) {
	query := `
		SELECT w.id, w.name, w.slug, w.description, w.owner_id, w.created_at, w.updated_at
		FROM workspaces w
		INNER JOIN workspace_members wm ON w.id = wm.workspace_id
		WHERE wm.user_id = $1 AND w.deleted_at IS NULL
		ORDER BY w.created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []model.Workspace
	for rows.Next() {
		var ws model.Workspace
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.Slug, &ws.Description, &ws.OwnerID, &ws.CreatedAt, &ws.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan workspace: %w", err)
		}
		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) Update(ctx context.Context, ws *model.Workspace) error {
	query := `
		UPDATE workspaces SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query, ws.ID, ws.Name, ws.Description).Scan(&ws.UpdatedAt)
}

func (r *WorkspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE workspaces SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *WorkspaceRepository) AddMember(ctx context.Context, member *model.WorkspaceMember) error {
	query := `
		INSERT INTO workspace_members (id, workspace_id, user_id, role)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`

	return r.db.QueryRow(ctx, query,
		member.ID, member.WorkspaceID, member.UserID, member.Role,
	).Scan(&member.CreatedAt)
}

func (r *WorkspaceRepository) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	query := `DELETE FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, workspaceID, userID)
	return err
}

func (r *WorkspaceRepository) UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role string) error {
	query := `UPDATE workspace_members SET role = $3 WHERE workspace_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, workspaceID, userID, role)
	return err
}

func (r *WorkspaceRepository) IsMember(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM workspace_members WHERE workspace_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, query, workspaceID, userID).Scan(&exists)
	return exists, err
}

func (r *WorkspaceRepository) GetMemberRole(ctx context.Context, workspaceID, userID uuid.UUID) (string, error) {
	var role string
	query := `SELECT role FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`
	err := r.db.QueryRow(ctx, query, workspaceID, userID).Scan(&role)
	return role, err
}
