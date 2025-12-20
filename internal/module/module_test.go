package module

import (
	"os"
	"path/filepath"
	"testing"
)

func createGoMod(t *testing.T, moduleLine string) string {
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(moduleLine), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.mod file: %v", err)
	}

	return tmpDir
}

func TestRead(t *testing.T) {
	tmpDir := createGoMod(t, "module github.com/example/project\n")

	moduleName, err := Read(tmpDir)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	expectedModuleName := "github.com/example/project"
	if moduleName != expectedModuleName {
		t.Errorf("Expected module name %q, got %q", expectedModuleName, moduleName)
	}
}

func TestReadNonExistentFile(t *testing.T) {
	_, err := Read("/non/existent/path")
	if err == nil {
		t.Fatal("Expected error for non-existent path, got nil")
	}
}

func TestReadMalformedGoMod(t *testing.T) {
	tmpDir := createGoMod(t, "modul github.com/example/project\n")

	_, err := Read(tmpDir)
	if err == nil {
		t.Fatal("Expected error for malformed go.mod, got nil")
	}
}
