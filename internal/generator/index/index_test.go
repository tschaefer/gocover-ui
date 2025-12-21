/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package index

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tschaefer/cover-ui/internal/coverage"
)

func TestAssets(t *testing.T) {
	outDir := t.TempDir()

	err := Assets(outDir)
	if err != nil {
		t.Fatalf("Assets() error = %v", err)
	}

	cssPath := filepath.Join(outDir, "style.css")
	jsPath := filepath.Join(outDir, "script.js")

	if _, err := os.Stat(cssPath); os.IsNotExist(err) {
		t.Errorf("Expected CSS file to exist at %s", cssPath)
	}

	if _, err := os.Stat(jsPath); os.IsNotExist(err) {
		t.Errorf("Expected JS file to exist at %s", jsPath)
	}
}

func TestGenerate(t *testing.T) {
	outDir := t.TempDir()

	files := []*coverage.FileMetrics{
		{FileName: "file1.go"},
		{FileName: "file2.go"},
	}
	module := "github.com/example/project"

	err := Generate(files, outDir, module)
	if err != nil {
		t.Fatalf("GenerateIndex() error = %v", err)
	}

	indexPath := filepath.Join(outDir, "index.html")

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Errorf("Expected index.html file to exist at %s", indexPath)
	}
}
