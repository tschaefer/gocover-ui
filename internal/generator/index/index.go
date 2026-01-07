/*
Copyright (c) Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package index

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/tschaefer/cover-ui/internal/coverage"
	"github.com/tschaefer/cover-ui/internal/generator/base"
	"github.com/tschaefer/cover-ui/internal/tree"
)

//go:embed assets/index.html
var indexHTML string

//go:embed assets/index.css
var indexCSS string

//go:embed assets/index.js
var indexJS string

// Generate creates the index page
func Generate(files []*coverage.FileMetrics, outDir string, module string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fileTree := tree.Build(files)

	metaJSON, err := json.Marshal(files)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}

	treeJSON, err := json.Marshal(fileTree)
	if err != nil {
		return fmt.Errorf("failed to marshal file tree to JSON: %w", err)
	}

	data := struct {
		Files    []*coverage.FileMetrics
		MetaJSON template.JS
		TreeJSON template.JS
		Module   string
	}{
		Files:    files,
		MetaJSON: template.JS(metaJSON),
		TreeJSON: template.JS(treeJSON),
		Module:   module,
	}

	return writeHTMLFile(outDir, data)
}

// Assets writes the css and javascript files to the output directory
func Assets(outDir string) error {
	mergedCSS := base.CSS + "\n\n" + indexCSS
	cssPath := filepath.Join(outDir, "style.css")
	if err := os.WriteFile(cssPath, []byte(mergedCSS), 0o644); err != nil {
		return fmt.Errorf("failed to write CSS file: %w", err)
	}

	jsPath := filepath.Join(outDir, "script.js")
	if err := os.WriteFile(jsPath, []byte(indexJS), 0o644); err != nil {
		return fmt.Errorf("failed to write JS file: %w", err)
	}

	return nil
}

// writeHTMLFile writes the index.html file to the output directory
func writeHTMLFile(outDir string, data any) error {
	tpl, err := template.New("base").Parse(base.HTML)
	if err != nil {
		return err
	}
	tpl, err = tpl.Parse(indexHTML)
	if err != nil {
		return err
	}

	outPath := filepath.Join(outDir, "index.html")
	w, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = w.Close()
	}()

	return tpl.Execute(w, data)
}
