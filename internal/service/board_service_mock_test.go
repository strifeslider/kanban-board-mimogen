package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/mock"
	"github.com/user/kanban-saas/pkg/model"
)

func newTestBoardService() (*BoardService, *mock.MockBoardRepo, *mock.MockColumnRepo, *mock.MockWorkspaceRepo) {
	boardRepo := mock.NewMockBoardRepo()
	columnRepo := mock.NewMockColumnRepo()
	workspaceRepo := mock.NewMockWorkspaceRepo()
	svc := NewBoardService(boardRepo, columnRepo, workspaceRepo)
	return svc, boardRepo, columnRepo, workspaceRepo
}

func TestBoardService_Create(t *testing.T) {
	svc, _, _, _ := newTestBoardService()
	ctx := context.Background()

	board, err := svc.Create(ctx, uuid.New(), uuid.New(), model.CreateBoardRequest{
		Name: "My Board",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if board.Name != "My Board" {
		t.Errorf("expected name 'My Board', got '%s'", board.Name)
	}
}

func TestBoardService_GetByID(t *testing.T) {
	svc, boardRepo, _, _ := newTestBoardService()
	ctx := context.Background()

	boardID := uuid.New()
	boardRepo.Boards[boardID] = &model.Board{
		ID:   boardID,
		Name: "Test",
	}

	board, err := svc.GetByID(ctx, boardID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if board.Name != "Test" {
		t.Errorf("expected name 'Test', got '%s'", board.Name)
	}
}

func TestBoardService_ListByWorkspace(t *testing.T) {
	svc, boardRepo, _, _ := newTestBoardService()
	ctx := context.Background()

	workspaceID := uuid.New()
	boardRepo.Boards[uuid.New()] = &model.Board{WorkspaceID: workspaceID, Name: "B1"}
	boardRepo.Boards[uuid.New()] = &model.Board{WorkspaceID: workspaceID, Name: "B2"}

	boards, err := svc.ListByWorkspace(ctx, workspaceID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(boards) != 2 {
		t.Errorf("expected 2 boards, got %d", len(boards))
	}
}

func TestBoardService_Update(t *testing.T) {
	svc, boardRepo, _, _ := newTestBoardService()
	ctx := context.Background()

	boardID := uuid.New()
	boardRepo.Boards[boardID] = &model.Board{
		ID:   boardID,
		Name: "Old",
	}

	newName := "New"
	board, err := svc.Update(ctx, boardID, model.UpdateBoardRequest{Name: &newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if board.Name != "New" {
		t.Errorf("expected name 'New', got '%s'", board.Name)
	}
}

func TestBoardService_Delete(t *testing.T) {
	svc, boardRepo, _, _ := newTestBoardService()
	ctx := context.Background()

	boardID := uuid.New()
	boardRepo.Boards[boardID] = &model.Board{ID: boardID}

	err := svc.Delete(ctx, boardID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBoardService_CreateColumn(t *testing.T) {
	svc, _, _, _ := newTestBoardService()
	ctx := context.Background()

	col, err := svc.CreateColumn(ctx, uuid.New(), model.CreateColumnRequest{
		Name: "TODO",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if col.Name != "TODO" {
		t.Errorf("expected name 'TODO', got '%s'", col.Name)
	}
	if col.Position != 0 {
		t.Errorf("expected position 0, got %d", col.Position)
	}
}

func TestBoardService_CreateColumn_WithPosition(t *testing.T) {
	svc, _, _, _ := newTestBoardService()
	ctx := context.Background()

	pos := 5
	col, err := svc.CreateColumn(ctx, uuid.New(), model.CreateColumnRequest{
		Name:     "DONE",
		Position: &pos,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if col.Position != 5 {
		t.Errorf("expected position 5, got %d", col.Position)
	}
}

func TestBoardService_ListColumns(t *testing.T) {
	svc, _, columnRepo, _ := newTestBoardService()
	ctx := context.Background()

	boardID := uuid.New()
	columnRepo.Columns[uuid.New()] = &model.Column{BoardID: boardID, Name: "C1"}
	columnRepo.Columns[uuid.New()] = &model.Column{BoardID: boardID, Name: "C2"}

	columns, err := svc.ListColumns(ctx, boardID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(columns))
	}
}

func TestBoardService_UpdateColumn(t *testing.T) {
	svc, _, columnRepo, _ := newTestBoardService()
	ctx := context.Background()

	colID := uuid.New()
	columnRepo.Columns[colID] = &model.Column{
		ID:   colID,
		Name: "Old",
	}

	newName := "New"
	col, err := svc.UpdateColumn(ctx, colID, model.UpdateColumnRequest{Name: &newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if col.Name != "New" {
		t.Errorf("expected name 'New', got '%s'", col.Name)
	}
}

func TestBoardService_DeleteColumn(t *testing.T) {
	svc, _, columnRepo, _ := newTestBoardService()
	ctx := context.Background()

	colID := uuid.New()
	columnRepo.Columns[colID] = &model.Column{ID: colID}

	err := svc.DeleteColumn(ctx, colID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBoardService_ReorderColumns(t *testing.T) {
	svc, _, _, _ := newTestBoardService()
	ctx := context.Background()

	ids := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	err := svc.ReorderColumns(ctx, uuid.New(), ids)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewBoardService(t *testing.T) {
	svc, _, _, _ := newTestBoardService()
	if svc == nil {
		t.Error("expected non-nil service")
	}
}
