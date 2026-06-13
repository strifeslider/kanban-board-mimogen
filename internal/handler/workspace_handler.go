package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/user/kanban-saas/pkg/auth"
	apperr "github.com/user/kanban-saas/pkg/errors"
	"github.com/user/kanban-saas/pkg/model"
	"github.com/user/kanban-saas/services/board/internal/service"
)

type WorkspaceHandler struct {
	workspaceService *service.WorkspaceService
}

func NewWorkspaceHandler(workspaceService *service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceService: workspaceService}
}

func (h *WorkspaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	var req model.CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	if req.Name == "" {
		apperr.RespondError(w, apperr.BadRequest("name is required"))
		return
	}

	ws, err := h.workspaceService.Create(r.Context(), userID, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusCreated, ws)
}

func (h *WorkspaceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	workspaces, err := h.workspaceService.ListByUser(r.Context(), userID)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	ws, err := h.workspaceService.GetByID(r.Context(), id)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	var req model.UpdateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	ws, err := h.workspaceService.Update(r.Context(), id, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	if err := h.workspaceService.Delete(r.Context(), id); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "workspace deleted"})
}

func (h *WorkspaceHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	var req struct {
		UserID uuid.UUID `json:"user_id"`
		Role   string    `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	member, err := h.workspaceService.AddMember(r.Context(), workspaceID, req.UserID, req.Role)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusCreated, member)
}

func (h *WorkspaceHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid user id"))
		return
	}

	if err := h.workspaceService.RemoveMember(r.Context(), workspaceID, userID); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "member removed"})
}

func (h *WorkspaceHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid user id"))
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	if err := h.workspaceService.UpdateMemberRole(r.Context(), workspaceID, userID, req.Role); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "role updated"})
}
