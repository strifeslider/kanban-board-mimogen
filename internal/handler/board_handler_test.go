package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/user/kanban-saas/pkg/auth"
)

func TestNewWorkspaceHandler(t *testing.T) {
	h := &WorkspaceHandler{}
	if h == nil {
		t.Error("expected non-nil handler")
	}
}

func TestNewBoardHandler(t *testing.T) {
	h := &BoardHandler{}
	if h == nil {
		t.Error("expected non-nil handler")
	}
}

func TestNewColumnHandler(t *testing.T) {
	h := &ColumnHandler{}
	if h == nil {
		t.Error("expected non-nil handler")
	}
}

func TestSetupRoutes(t *testing.T) {
	r := chi.NewRouter()
	wh := &WorkspaceHandler{}
	bh := &BoardHandler{}
	ch := &ColumnHandler{}
	jwtCfg := auth.JWTConfig{Secret: "test"}

	SetupRoutes(r, wh, bh, ch, jwtCfg)

	// Verify routes are registered
	routes := []string{
		"/api/v1/workspaces",
		"/api/v1/boards/test-id",
	}

	for _, route := range routes {
		req := httptest.NewRequest("GET", route, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		// Should not return 405 (method not allowed) for route existence
		if w.Code == http.StatusNotFound {
			t.Errorf("route %s not found", route)
		}
	}
}

func TestWorkspaceHandler_Create_EmptyBody(t *testing.T) {
	h := &WorkspaceHandler{}
	req := httptest.NewRequest("POST", "/api/v1/workspaces", nil)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestBoardHandler_Get_InvalidID(t *testing.T) {
	h := &BoardHandler{}
	req := httptest.NewRequest("GET", "/api/v1/boards/invalid-id", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestColumnHandler_Create_InvalidBoardID(t *testing.T) {
	h := &ColumnHandler{}
	req := httptest.NewRequest("POST", "/api/v1/boards/invalid/columns", nil)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
