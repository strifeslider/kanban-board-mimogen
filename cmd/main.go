package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/auth"
	"github.com/user/kanban-saas/pkg/database"
	appmiddleware "github.com/user/kanban-saas/pkg/middleware"
	ws "github.com/user/kanban-saas/pkg/websocket"
	"github.com/user/kanban-saas/services/board/internal/handler"
	"github.com/user/kanban-saas/services/board/internal/repository"
	"github.com/user/kanban-saas/services/board/internal/service"
)

func main() {
	env := getEnv("ENV", "local")
	port := getEnv("PORT", "8082")

	logger := setupLogger(env)
	logger.Info("starting board service", "env", env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewPostgresPool(ctx, database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "kanban"),
		Password: getEnv("DB_PASSWORD", "kanban_dev_password"),
		Database: getEnv("DB_NAME", "kanban_board"),
		MaxConns: 10,
		MinConns: 2,
	})
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	runMigrations(ctx, db, logger)

	jwtCfg := auth.JWTConfig{
		Secret: getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
	}

	hub := ws.NewHub(logger)
	go hub.Run()

	workspaceRepo := repository.NewWorkspaceRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	columnRepo := repository.NewColumnRepository(db)

	workspaceService := service.NewWorkspaceService(workspaceRepo)
	boardService := service.NewBoardService(boardRepo, columnRepo, workspaceRepo)

	workspaceHandler := handler.NewWorkspaceHandler(workspaceService)
	boardHandler := handler.NewBoardHandler(boardService)
	columnHandler := handler.NewColumnHandler(boardService)

	r := chi.NewRouter()

	allowedOrigins := appmiddleware.ParseOrigins(getEnv("CORS_ORIGINS", "http://localhost:3000"))
	r.Use(appmiddleware.CORS(allowedOrigins))
	r.Use(appmiddleware.Logging(logger))
	r.Use(appmiddleware.Recovery(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r, jwtCfg, logger)
	})

	handler.SetupRoutes(r, workspaceHandler, boardHandler, columnHandler, jwtCfg)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("board service listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down board service...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}
	logger.Info("board service stopped")
}

func handleWebSocket(hub *ws.Hub, w http.ResponseWriter, r *http.Request, jwtCfg auth.JWTConfig, logger *slog.Logger) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	_, err := auth.ValidateToken(jwtCfg, token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	upgrader := ws.NewUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("websocket upgrade failed", "error", err)
		return
	}

	client := ws.NewClient(hub, conn)
	hub.Register(client, "global")

	go client.WritePump()
	go client.ReadPump()
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler
	switch env {
	case "local":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "dev":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	return slog.New(handler)
}

func runMigrations(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger) {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS workspaces (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			owner_id UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);`,
		`CREATE INDEX IF NOT EXISTS idx_workspaces_owner ON workspaces(owner_id) WHERE deleted_at IS NULL;`,
		`CREATE TABLE IF NOT EXISTS workspace_members (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			user_id UUID NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'member',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(workspace_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS boards (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created_by UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);`,
		`CREATE INDEX IF NOT EXISTS idx_boards_workspace ON boards(workspace_id) WHERE deleted_at IS NULL;`,
		`CREATE TABLE IF NOT EXISTS board_members (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			board_id UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
			user_id UUID NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'member',
			UNIQUE(board_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS columns (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			board_id UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			position INTEGER NOT NULL DEFAULT 0,
			color VARCHAR(7),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_columns_board ON columns(board_id);`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(ctx, m); err != nil {
			logger.Error("migration failed", "error", err)
			os.Exit(1)
		}
	}
	logger.Info("migrations completed")
}
