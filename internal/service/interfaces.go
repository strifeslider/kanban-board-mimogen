package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/model"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, ws *model.Workspace) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Workspace, error)
	Update(ctx context.Context, ws *model.Workspace) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddMember(ctx context.Context, member *model.WorkspaceMember) error
	RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role string) error
	IsMember(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error)
	GetMemberRole(ctx context.Context, workspaceID, userID uuid.UUID) (string, error)
}

type BoardRepository interface {
	Create(ctx context.Context, board *model.Board) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Board, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.Board, error)
	Update(ctx context.Context, board *model.Board) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddMember(ctx context.Context, member *model.BoardMember) error
	RemoveMember(ctx context.Context, boardID, userID uuid.UUID) error
	IsMember(ctx context.Context, boardID, userID uuid.UUID) (bool, error)
}

type ColumnRepository interface {
	Create(ctx context.Context, col *model.Column) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Column, error)
	ListByBoard(ctx context.Context, boardID uuid.UUID) ([]model.Column, error)
	Update(ctx context.Context, col *model.Column) error
	Delete(ctx context.Context, id uuid.UUID) error
	Reorder(ctx context.Context, columnIDs []uuid.UUID) error
	GetMaxPosition(ctx context.Context, boardID uuid.UUID) (int, error)
}
