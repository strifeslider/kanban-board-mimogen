package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/mock"
	"github.com/user/kanban-saas/pkg/model"
)

func newTestWorkspaceService() (*WorkspaceService, *mock.MockWorkspaceRepo) {
	workspaceRepo := mock.NewMockWorkspaceRepo()
	svc := NewWorkspaceService(workspaceRepo)
	return svc, workspaceRepo
}

func TestWorkspaceService_Create(t *testing.T) {
	svc, _ := newTestWorkspaceService()
	ctx := context.Background()

	ws, err := svc.Create(ctx, uuid.New(), model.CreateWorkspaceRequest{
		Name: "My Workspace",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Name != "My Workspace" {
		t.Errorf("expected name 'My Workspace', got '%s'", ws.Name)
	}
	if ws.Slug != "my-workspace" {
		t.Errorf("expected slug 'my-workspace', got '%s'", ws.Slug)
	}
}

func TestWorkspaceService_GetByID(t *testing.T) {
	svc, workspaceRepo := newTestWorkspaceService()
	ctx := context.Background()

	wsID := uuid.New()
	workspaceRepo.Workspaces[wsID] = &model.Workspace{
		ID:   wsID,
		Name: "Test",
	}

	ws, err := svc.GetByID(ctx, wsID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Name != "Test" {
		t.Errorf("expected name 'Test', got '%s'", ws.Name)
	}
}

func TestWorkspaceService_ListByUser(t *testing.T) {
	svc, workspaceRepo := newTestWorkspaceService()
	ctx := context.Background()

	userID := uuid.New()
	workspaceRepo.Workspaces[uuid.New()] = &model.Workspace{OwnerID: userID, Name: "WS1"}
	workspaceRepo.Workspaces[uuid.New()] = &model.Workspace{OwnerID: userID, Name: "WS2"}

	wsList, err := svc.ListByUser(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(wsList) != 2 {
		t.Errorf("expected 2 workspaces, got %d", len(wsList))
	}
}

func TestWorkspaceService_Update(t *testing.T) {
	svc, workspaceRepo := newTestWorkspaceService()
	ctx := context.Background()

	wsID := uuid.New()
	workspaceRepo.Workspaces[wsID] = &model.Workspace{
		ID:   wsID,
		Name: "Old",
	}

	newName := "New"
	ws, err := svc.Update(ctx, wsID, model.UpdateWorkspaceRequest{Name: &newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Name != "New" {
		t.Errorf("expected name 'New', got '%s'", ws.Name)
	}
}

func TestWorkspaceService_Delete(t *testing.T) {
	svc, workspaceRepo := newTestWorkspaceService()
	ctx := context.Background()

	wsID := uuid.New()
	workspaceRepo.Workspaces[wsID] = &model.Workspace{ID: wsID}

	err := svc.Delete(ctx, wsID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceService_AddMember(t *testing.T) {
	svc, _ := newTestWorkspaceService()
	ctx := context.Background()

	member, err := svc.AddMember(ctx, uuid.New(), uuid.New(), "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if member.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", member.Role)
	}
}

func TestWorkspaceService_RemoveMember(t *testing.T) {
	svc, _ := newTestWorkspaceService()
	ctx := context.Background()

	err := svc.RemoveMember(ctx, uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceService_UpdateMemberRole(t *testing.T) {
	svc, _ := newTestWorkspaceService()
	ctx := context.Background()

	err := svc.UpdateMemberRole(ctx, uuid.New(), uuid.New(), "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceService_IsMember(t *testing.T) {
	svc, workspaceRepo := newTestWorkspaceService()
	ctx := context.Background()

	wsID := uuid.New()
	userID := uuid.New()
	workspaceRepo.Members[wsID] = []model.WorkspaceMember{
		{WorkspaceID: wsID, UserID: userID, Role: "member"},
	}

	isMember, err := svc.IsMember(ctx, wsID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isMember {
		t.Error("expected user to be member")
	}
}

func TestWorkspaceService_GetMemberRole(t *testing.T) {
	svc, workspaceRepo := newTestWorkspaceService()
	ctx := context.Background()

	wsID := uuid.New()
	userID := uuid.New()
	workspaceRepo.Members[wsID] = []model.WorkspaceMember{
		{WorkspaceID: wsID, UserID: userID, Role: "admin"},
	}

	role, err := svc.GetMemberRole(ctx, wsID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", role)
	}
}

func TestNewWorkspaceService(t *testing.T) {
	svc, _ := newTestWorkspaceService()
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestWorkspaceService_SlugGeneration(t *testing.T) {
	svc, _ := newTestWorkspaceService()
	ctx := context.Background()

	ws, _ := svc.Create(ctx, uuid.New(), model.CreateWorkspaceRequest{Name: "Hello World!"})
	if ws.Slug != "hello-world" {
		t.Errorf("expected slug 'hello-world', got '%s'", ws.Slug)
	}
}

func TestWorkspaceService_EmptyName(t *testing.T) {
	svc, _ := newTestWorkspaceService()
	ctx := context.Background()

	ws, _ := svc.Create(ctx, uuid.New(), model.CreateWorkspaceRequest{Name: ""})
	if ws.Slug == "" {
		t.Error("expected non-empty slug for empty name")
	}
}
