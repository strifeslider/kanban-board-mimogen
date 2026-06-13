package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/user/kanban-saas/pkg/model"
)

type WorkspaceService struct {
	workspaceRepo WorkspaceRepository
}

func NewWorkspaceService(workspaceRepo WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{workspaceRepo: workspaceRepo}
}

func (s *WorkspaceService) Create(ctx context.Context, userID uuid.UUID, req model.CreateWorkspaceRequest) (*model.Workspace, error) {
	slug := generateSlug(req.Name)

	ws := &model.Workspace{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		OwnerID:     userID,
	}

	if err := s.workspaceRepo.Create(ctx, ws); err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}

	member := &model.WorkspaceMember{
		ID:          uuid.New(),
		WorkspaceID: ws.ID,
		UserID:      userID,
		Role:        "owner",
	}
	if err := s.workspaceRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("add owner member: %w", err)
	}

	return ws, nil
}

func (s *WorkspaceService) GetByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error) {
	return s.workspaceRepo.GetByID(ctx, id)
}

func (s *WorkspaceService) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Workspace, error) {
	return s.workspaceRepo.ListByUser(ctx, userID)
}

func (s *WorkspaceService) Update(ctx context.Context, id uuid.UUID, req model.UpdateWorkspaceRequest) (*model.Workspace, error) {
	ws, err := s.workspaceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		ws.Name = *req.Name
	}
	if req.Description != nil {
		ws.Description = req.Description
	}

	if err := s.workspaceRepo.Update(ctx, ws); err != nil {
		return nil, fmt.Errorf("update workspace: %w", err)
	}

	return ws, nil
}

func (s *WorkspaceService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.workspaceRepo.Delete(ctx, id)
}

func (s *WorkspaceService) AddMember(ctx context.Context, workspaceID, userID uuid.UUID, role string) (*model.WorkspaceMember, error) {
	member := &model.WorkspaceMember{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
	}

	if err := s.workspaceRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("add member: %w", err)
	}

	return member, nil
}

func (s *WorkspaceService) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	return s.workspaceRepo.RemoveMember(ctx, workspaceID, userID)
}

func (s *WorkspaceService) UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role string) error {
	return s.workspaceRepo.UpdateMemberRole(ctx, workspaceID, userID, role)
}

func (s *WorkspaceService) IsMember(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error) {
	return s.workspaceRepo.IsMember(ctx, workspaceID, userID)
}

func (s *WorkspaceService) GetMemberRole(ctx context.Context, workspaceID, userID uuid.UUID) (string, error) {
	return s.workspaceRepo.GetMemberRole(ctx, workspaceID, userID)
}

func generateSlug(name string) string {
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug := strings.ToLower(name)
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = uuid.New().String()[:8]
	}
	return slug
}
