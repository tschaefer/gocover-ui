package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tschaefer/cover-ui/internal/coverage"
)

func TestAssets(t *testing.T) {
	outputDir := t.TempDir()
	err := Assets(outputDir)
	if err != nil {
		t.Fatalf("Assets() returned an error: %v", err)
	}

	cssPath := filepath.Join(outputDir, "style.css")
	if _, err := os.Stat(cssPath); os.IsNotExist(err) {
		t.Errorf("Expected CSS file does not exist: %s", cssPath)
	}

	jsPath := filepath.Join(outputDir, "script.js")
	if _, err := os.Stat(jsPath); os.IsNotExist(err) {
		t.Errorf("Expected JS file does not exist: %s", jsPath)
	}
}

func TestGenerate(t *testing.T) {
	outputDir := t.TempDir()
	filesDir := filepath.Join(outputDir, "files")

	sourceFile, err := os.CreateTemp("", "source-*.go")
	if err != nil {
		t.Fatalf("Failed to create temp source file: %v", err)
	}
	defer func() {
		_ = os.Remove(sourceFile.Name())
	}()

	sourceContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	if _, err := sourceFile.WriteString(sourceContent); err != nil {
		t.Fatalf("Failed to write to temp source file: %v", err)
	}
	_ = sourceFile.Close()

	fileMetrics := &coverage.FileMetrics{
		LocalPath: sourceFile.Name(),
	}

	err = Generate(fileMetrics, filesDir)
	if err != nil {
		t.Fatalf("Generate() returned an error: %v", err)
	}

	generatedFilePath := filepath.Join(filesDir, "tmp", strings.ReplaceAll(filepath.Base(sourceFile.Name()), ".go", ".html"))
	if _, err := os.Stat(generatedFilePath); os.IsNotExist(err) {
		t.Errorf("Expected generated HTML file does not exist: %s", generatedFilePath)
	}
}
