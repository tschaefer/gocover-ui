package coverage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/cover"
)

func createSourceFile(t *testing.T) (string, string) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "test.go")
	sourceContent := `package main

func main() {
	println("hello")
}
`
	if err := os.WriteFile(sourceFile, []byte(sourceContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return sourceFile, tmpDir
}

func TestAnalyze(t *testing.T) {
	sourceFile, tmpDir := createSourceFile(t)

	profile := &cover.Profile{
		FileName: sourceFile,
		Mode:     "set",
		Blocks: []cover.ProfileBlock{
			{
				StartLine: 3,
				StartCol:  13,
				EndLine:   5,
				EndCol:    2,
				NumStmt:   2,
				Count:     1,
			},
		},
	}

	metrics, err := Analyze(profile, "/", tmpDir)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Expected non-nil metrics")
	}
	if metrics.FileName != sourceFile {
		t.Errorf("Expected FileName %s, got %s", sourceFile, metrics.FileName)
	}
	if metrics.LocalPath != filepath.Join(tmpDir, "test.go") {
		t.Errorf("Expected LocalPath test.go, got %s", metrics.LocalPath)
	}
	if metrics.TrackedLines != 3 {
		t.Errorf("Expected TrackedLines 3, got %d", metrics.TrackedLines)
	}
	if metrics.CoveredLines != 3 {
		t.Errorf("Expected CoveredLines 3, got %d", metrics.CoveredLines)
	}
	if metrics.PartialLines != 0 {
		t.Errorf("Expected PartialLines 0, got %d", metrics.PartialLines)
	}
	if metrics.MissedLines != 0 {
		t.Errorf("Expected MissedLines 0, got %d", metrics.MissedLines)
	}
	if metrics.TotalStmts != 2 {
		t.Errorf("Expected TotalStmts 2, got %d", metrics.TotalStmts)
	}
	if metrics.CoveredStmts != 2 {
		t.Errorf("Expected CoveredStmts 2, got %d", metrics.CoveredStmts)
	}
	if metrics.CoveragePct != 100.0 {
		t.Errorf("Expected CoveragePct 100.0, got %.2f", metrics.CoveragePct)
	}
}

func TestModuleMismatch(t *testing.T) {
	sourceFile, tmpDir := createSourceFile(t)

	profile := &cover.Profile{
		FileName: sourceFile,
		Mode:     "set",
		Blocks:   []cover.ProfileBlock{},
	}

	_, err := Analyze(profile, "/different/module", tmpDir)
	if err == nil {
		t.Fatal("Expected error due to module mismatch, got nil")
	}
	if err.Error() != "failed to match module /different/module" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestSourceReadError(t *testing.T) {
	profile := &cover.Profile{
		FileName: "/non/existent/file.go",
		Mode:     "set",
		Blocks:   []cover.ProfileBlock{},
	}

	_, err := Analyze(profile, "/", ".")
	if err == nil {
		t.Fatal("Expected error due to file read failure, got nil")
	}
	expectedPrefix := "failed to read source file /non/existent/file.go"
	if !strings.HasPrefix(err.Error(), expectedPrefix) {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestLineStatus(t *testing.T) {
	tests := []struct {
		name   string
		status LineStatus
		want   int
	}{
		{"missed", Missed, 0},
		{"partial", Partial, 1},
		{"covered", Covered, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.status) != tt.want {
				t.Errorf("Expected %s to be %d, got %d", tt.name, tt.want, int(tt.status))
			}
		})
	}
}

func TestLineStatusString(t *testing.T) {
	tests := []struct {
		status LineStatus
		want   string
	}{
		{Missed, "missed"},
		{Partial, "partial"},
		{Covered, "covered"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got, err := tt.status.String()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Expected %s, got %s", tt.want, got)
			}
		})
	}

	var ls LineStatus = 99
	_, err := ls.String()
	if err == nil {
		t.Fatal("Expected error for invalid LineStatus, got nil")
	}
}
