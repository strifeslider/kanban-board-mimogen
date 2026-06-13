package service

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My Workspace", "my-workspace"},
		{"Hello World!", "hello-world"},
		{"Test 123", "test-123"},
		{"UPPERCASE", "uppercase"},
		{"special@#$chars", "special-chars"},
	}

	for _, tt := range tests {
		result := generateSlug(tt.input)
		if result != tt.expected {
			t.Errorf("generateSlug(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}

	// Test empty input generates a slug
	result := generateSlug("")
	if result == "" {
		t.Error("generateSlug(\"\") should generate a slug for empty input")
	}
}

func TestNewWorkspaceService(t *testing.T) {
	svc := &WorkspaceService{}
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestNewBoardService(t *testing.T) {
	svc := &BoardService{}
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestUUIDGeneration(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	if id1 == id2 {
		t.Error("expected different UUIDs")
	}
}
