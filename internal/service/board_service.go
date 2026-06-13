package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/user/kanban-saas/pkg/model"
)

type BoardService struct {
	boardRepo    BoardRepository
	columnRepo   ColumnRepository
	workspaceRepo WorkspaceRepository
}

func NewBoardService(
	boardRepo BoardRepository,
	columnRepo ColumnRepository,
	workspaceRepo WorkspaceRepository,
) *BoardService {
	return &BoardService{
		boardRepo:     boardRepo,
		columnRepo:    columnRepo,
		workspaceRepo: workspaceRepo,
	}
}

func (s *BoardService) Create(ctx context.Context, userID, workspaceID uuid.UUID, req model.CreateBoardRequest) (*model.Board, error) {
	board := &model.Board{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	}

	if err := s.boardRepo.Create(ctx, board); err != nil {
		return nil, fmt.Errorf("create board: %w", err)
	}

	member := &model.BoardMember{
		ID:      uuid.New(),
		BoardID: board.ID,
		UserID:  userID,
		Role:    "admin",
	}
	if err := s.boardRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("add board member: %w", err)
	}

	return board, nil
}

func (s *BoardService) GetByID(ctx context.Context, id uuid.UUID) (*model.Board, error) {
	return s.boardRepo.GetByID(ctx, id)
}

func (s *BoardService) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.Board, error) {
	return s.boardRepo.ListByWorkspace(ctx, workspaceID)
}

func (s *BoardService) Update(ctx context.Context, id uuid.UUID, req model.UpdateBoardRequest) (*model.Board, error) {
	board, err := s.boardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		board.Name = *req.Name
	}
	if req.Description != nil {
		board.Description = req.Description
	}

	if err := s.boardRepo.Update(ctx, board); err != nil {
		return nil, fmt.Errorf("update board: %w", err)
	}

	return board, nil
}

func (s *BoardService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.boardRepo.Delete(ctx, id)
}

func (s *BoardService) CreateColumn(ctx context.Context, boardID uuid.UUID, req model.CreateColumnRequest) (*model.Column, error) {
	pos := 0
	if req.Position != nil {
		pos = *req.Position
	} else {
		maxPos, err := s.columnRepo.GetMaxPosition(ctx, boardID)
		if err != nil {
			return nil, fmt.Errorf("get max position: %w", err)
		}
		pos = maxPos + 1
	}

	col := &model.Column{
		ID:       uuid.New(),
		BoardID:  boardID,
		Name:     req.Name,
		Position: pos,
		Color:    req.Color,
	}

	if err := s.columnRepo.Create(ctx, col); err != nil {
		return nil, fmt.Errorf("create column: %w", err)
	}

	return col, nil
}

func (s *BoardService) GetColumn(ctx context.Context, id uuid.UUID) (*model.Column, error) {
	return s.columnRepo.GetByID(ctx, id)
}

func (s *BoardService) ListColumns(ctx context.Context, boardID uuid.UUID) ([]model.Column, error) {
	return s.columnRepo.ListByBoard(ctx, boardID)
}

func (s *BoardService) UpdateColumn(ctx context.Context, id uuid.UUID, req model.UpdateColumnRequest) (*model.Column, error) {
	col, err := s.columnRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		col.Name = *req.Name
	}
	if req.Position != nil {
		col.Position = *req.Position
	}
	if req.Color != nil {
		col.Color = req.Color
	}

	if err := s.columnRepo.Update(ctx, col); err != nil {
		return nil, fmt.Errorf("update column: %w", err)
	}

	return col, nil
}

func (s *BoardService) DeleteColumn(ctx context.Context, id uuid.UUID) error {
	return s.columnRepo.Delete(ctx, id)
}

func (s *BoardService) ReorderColumns(ctx context.Context, boardID uuid.UUID, columnIDs []uuid.UUID) error {
	return s.columnRepo.Reorder(ctx, columnIDs)
}
