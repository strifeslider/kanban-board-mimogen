package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/user/kanban-saas/pkg/auth"
)

func SetupRoutes(
	r chi.Router,
	workspaceHandler *WorkspaceHandler,
	boardHandler *BoardHandler,
	columnHandler *ColumnHandler,
	jwtCfg auth.JWTConfig,
) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(auth.RequireAuth(jwtCfg))

		r.Route("/workspaces", func(r chi.Router) {
			r.Post("/", workspaceHandler.Create)
			r.Get("/", workspaceHandler.List)
			r.Get("/{id}", workspaceHandler.Get)
			r.Put("/{id}", workspaceHandler.Update)
			r.Delete("/{id}", workspaceHandler.Delete)

			r.Post("/{id}/members", workspaceHandler.AddMember)
			r.Delete("/{id}/members/{userId}", workspaceHandler.RemoveMember)
			r.Put("/{id}/members/{userId}", workspaceHandler.UpdateMemberRole)

			r.Route("/{workspaceId}/boards", func(r chi.Router) {
				r.Post("/", boardHandler.Create)
				r.Get("/", boardHandler.List)
			})
		})

		r.Route("/boards", func(r chi.Router) {
			r.Get("/{id}", boardHandler.Get)
			r.Put("/{id}", boardHandler.Update)
			r.Delete("/{id}", boardHandler.Delete)

			r.Post("/{boardId}/columns", columnHandler.Create)
			r.Get("/{boardId}/columns", columnHandler.List)
			r.Put("/{boardId}/columns/reorder", columnHandler.Reorder)
		})

		r.Route("/columns", func(r chi.Router) {
			r.Put("/{id}", columnHandler.Update)
			r.Delete("/{id}", columnHandler.Delete)
		})
	})
}
