package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/user/kanban-saas/pkg/auth"
)

func TestWorkspaceHandler_List_NoAuth(t *testing.T) {
	h := &WorkspaceHandler{}
	req := httptest.NewRequest("GET", "/api/v1/workspaces", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestWorkspaceHandler_Get_InvalidID(t *testing.T) {
	h := &WorkspaceHandler{}
	req := httptest.NewRequest("GET", "/api/v1/workspaces/invalid", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestWorkspaceHandler_Update_InvalidID(t *testing.T) {
	h := &WorkspaceHandler{}
	body := map[string]string{"name": "New Name"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/workspaces/invalid", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestWorkspaceHandler_Delete_InvalidID(t *testing.T) {
	h := &WorkspaceHandler{}
	req := httptest.NewRequest("DELETE", "/api/v1/workspaces/invalid", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestWorkspaceHandler_AddMember_InvalidID(t *testing.T) {
	h := &WorkspaceHandler{}
	body := map[string]string{"user_id": "invalid"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/workspaces/invalid/members", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	h.AddMember(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestBoardHandler_List_InvalidWorkspaceID(t *testing.T) {
	h := &BoardHandler{}
	req := httptest.NewRequest("GET", "/api/v1/workspaces/invalid/boards", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestBoardHandler_Update_InvalidID(t *testing.T) {
	h := &BoardHandler{}
	body := map[string]string{"name": "New Name"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/boards/invalid", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestBoardHandler_Delete_InvalidID(t *testing.T) {
	h := &BoardHandler{}
	req := httptest.NewRequest("DELETE", "/api/v1/boards/invalid", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestColumnHandler_List_InvalidBoardID(t *testing.T) {
	h := &ColumnHandler{}
	req := httptest.NewRequest("GET", "/api/v1/boards/invalid/columns", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestColumnHandler_Update_InvalidID(t *testing.T) {
	h := &ColumnHandler{}
	body := map[string]string{"name": "New Name"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/columns/invalid", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestColumnHandler_Delete_InvalidID(t *testing.T) {
	h := &ColumnHandler{}
	req := httptest.NewRequest("DELETE", "/api/v1/columns/invalid", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestColumnHandler_Reorder_InvalidBoardID(t *testing.T) {
	h := &ColumnHandler{}
	body := map[string][]string{"column_ids": {"invalid"}}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/boards/invalid/columns/reorder", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	h.Reorder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSetupRoutes_Complete(t *testing.T) {
	r := chi.NewRouter()
	wh := &WorkspaceHandler{}
	bh := &BoardHandler{}
	ch := &ColumnHandler{}
	jwtCfg := auth.JWTConfig{Secret: "test"}

	SetupRoutes(r, wh, bh, ch, jwtCfg)

	routes := []string{
		"/api/v1/workspaces",
		"/api/v1/workspaces/test-id",
		"/api/v1/workspaces/test-id/boards",
		"/api/v1/boards/test-id",
		"/api/v1/boards/test-id/columns",
		"/api/v1/columns/test-id",
	}

	for _, route := range routes {
		req := httptest.NewRequest("GET", route, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code == http.StatusNotFound {
			t.Errorf("route %s not found", route)
		}
	}
}
