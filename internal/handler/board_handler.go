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

type BoardHandler struct {
	boardService *service.BoardService
}

func NewBoardHandler(boardService *service.BoardService) *BoardHandler {
	return &BoardHandler{boardService: boardService}
}

func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	workspaceID, err := uuid.Parse(chi.URLParam(r, "workspaceId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	var req model.CreateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	board, err := h.boardService.Create(r.Context(), userID, workspaceID, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusCreated, board)
}

func (h *BoardHandler) List(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "workspaceId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	boards, err := h.boardService.ListByWorkspace(r.Context(), workspaceID)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, boards)
}

func (h *BoardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid board id"))
		return
	}

	board, err := h.boardService.GetByID(r.Context(), id)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, board)
}

func (h *BoardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid board id"))
		return
	}

	var req model.UpdateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	board, err := h.boardService.Update(r.Context(), id, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, board)
}

func (h *BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid board id"))
		return
	}

	if err := h.boardService.Delete(r.Context(), id); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "board deleted"})
}
