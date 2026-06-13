package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	apperr "github.com/user/kanban-saas/pkg/errors"
	"github.com/user/kanban-saas/pkg/model"
	"github.com/user/kanban-saas/services/board/internal/service"
)

type ColumnHandler struct {
	boardService *service.BoardService
}

func NewColumnHandler(boardService *service.BoardService) *ColumnHandler {
	return &ColumnHandler{boardService: boardService}
}

func (h *ColumnHandler) Create(w http.ResponseWriter, r *http.Request) {
	boardID, err := uuid.Parse(chi.URLParam(r, "boardId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid board id"))
		return
	}

	var req model.CreateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	col, err := h.boardService.CreateColumn(r.Context(), boardID, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusCreated, col)
}

func (h *ColumnHandler) List(w http.ResponseWriter, r *http.Request) {
	boardID, err := uuid.Parse(chi.URLParam(r, "boardId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid board id"))
		return
	}

	columns, err := h.boardService.ListColumns(r.Context(), boardID)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, columns)
}

func (h *ColumnHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid column id"))
		return
	}

	var req model.UpdateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	col, err := h.boardService.UpdateColumn(r.Context(), id, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, col)
}

func (h *ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid column id"))
		return
	}

	if err := h.boardService.DeleteColumn(r.Context(), id); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "column deleted"})
}

func (h *ColumnHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	boardID, err := uuid.Parse(chi.URLParam(r, "boardId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid board id"))
		return
	}

	var req model.ReorderColumnsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	if err := h.boardService.ReorderColumns(r.Context(), boardID, req.ColumnIDs); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "columns reordered"})
}
