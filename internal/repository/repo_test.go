package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/model"
)

func TestWorkspaceRepository_New(t *testing.T) {
	repo := &WorkspaceRepository{}
	if repo == nil {
		t.Error("expected non-nil repo")
	}
}

func TestBoardRepository_New(t *testing.T) {
	repo := &BoardRepository{}
	if repo == nil {
		t.Error("expected non-nil repo")
	}
}

func TestColumnRepository_New(t *testing.T) {
	repo := &ColumnRepository{}
	if repo == nil {
		t.Error("expected non-nil repo")
	}
}

func TestWorkspaceRepository_Model(t *testing.T) {
	ws := &model.Workspace{
		ID:      uuid.New(),
		Name:    "Test",
		Slug:    "test",
		OwnerID: uuid.New(),
	}
	if ws.Name != "Test" {
		t.Error("name mismatch")
	}
}

func TestBoardRepository_Model(t *testing.T) {
	board := &model.Board{
		ID:          uuid.New(),
		WorkspaceID: uuid.New(),
		Name:        "Board",
		CreatedBy:   uuid.New(),
	}
	if board.Name != "Board" {
		t.Error("name mismatch")
	}
}

func TestColumnRepository_Model(t *testing.T) {
	col := &model.Column{
		ID:       uuid.New(),
		BoardID:  uuid.New(),
		Name:     "TODO",
		Position: 0,
	}
	if col.Name != "TODO" {
		t.Error("name mismatch")
	}
}
