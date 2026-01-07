/*
Copyright (c) Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package file

import (
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/tschaefer/cover-ui/internal/coverage"
	"github.com/tschaefer/cover-ui/internal/generator/base"
)

//go:embed assets/file.html
var fileHTML string

//go:embed assets/file.css
var fileCSS string

//go:embed assets/file.js
var fileJS string

// Generate creates a file detail page
func Generate(f *coverage.FileMetrics, filesDir string) error {
	if err := os.MkdirAll(filesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create files directory: %w", err)
	}

	source, err := os.ReadFile(f.LocalPath)
	if err != nil {
		return fmt.Errorf("failed to read source file %q: %w", f.LocalPath, err)
	}

	data := struct {
		File  *coverage.FileMetrics
		Lines []string
	}{
		File:  f,
		Lines: getLines(source),
	}

	if err := writeHTMLFile(data, filesDir, f); err != nil {
		return fmt.Errorf("failed to write file detail page for %q: %w", f.LocalPath, err)
	}

	return nil
}

// Assets writes the css and javascript files to the output directory
func Assets(filesDir string) error {
	if err := os.MkdirAll(filesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create files directory: %w", err)
	}

	mergedCSS := base.CSS + "\n\n" + fileCSS
	cssPath := filepath.Join(filesDir, "style.css")
	if err := os.WriteFile(cssPath, []byte(mergedCSS), 0o644); err != nil {
		return fmt.Errorf("failed to write css file: %w", err)
	}

	jsPath := filepath.Join(filesDir, "script.js")
	if err := os.WriteFile(jsPath, []byte(fileJS), 0o644); err != nil {
		return fmt.Errorf("failed to write javascript file: %w", err)
	}

	return nil
}

// getLines splits the source into lines and removes trailing empty ones
func getLines(source []byte) []string {
	lines := strings.Split(string(source), "\n")

	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}

// writeHTMLFile writes the detail HTML page
func writeHTMLFile(data any, filesDir string, f *coverage.FileMetrics) error {
	tpl, err := template.New("base").Funcs(template.FuncMap{
		"escape":     __EscapeSourceLine,
		"lineClass":  __AddSourceLineClass,
		"lineMarker": __AddLineMarker,
		"inc":        __IncByOne,
		"indexPath":  func() string { return __GetRelativePath(f.LocalPath, "../index.html") },
		"cssPath":    func() string { return __GetRelativePath(f.LocalPath, "style.css") },
		"scriptPath": func() string { return __GetRelativePath(f.LocalPath, "script.js") },
	}).Parse(base.HTML)
	if err != nil {
		return err
	}
	tpl, err = tpl.Parse(fileHTML)
	if err != nil {
		return err
	}

	outPath := filepath.Join(filesDir, f.LocalPath)
	outPath = strings.TrimSuffix(outPath, filepath.Ext(outPath)) + ".html"

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	w, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = w.Close()
	}()

	return tpl.Execute(w, data)
}

// Template helper functions

// __EscapesSourceLines escapes source line for HTML output
func __EscapeSourceLine(str string) template.HTML {
	if str == "" {
		return template.HTML("&nbsp;")
	}

	return template.HTML(html.EscapeString(str))
}

// __AddSourceLineClass adds a CSS class based on the coverage status of a
// source code line
func __AddSourceLineClass(idx int, statuses []int) string {
	if idx < 1 || idx >= len(statuses) {
		return "not-tracked"
	}

	str, err := coverage.LineStatus(statuses[idx]).String()
	if err != nil {
		return "not-tracked"
	}

	return str
}

// __AddLineMarker adds exclamation mark markers for partial and missed lines
func __AddLineMarker(idx int, statuses []int) string {
	if idx < 1 || idx >= len(statuses) {
		return ""
	}

	switch coverage.LineStatus(statuses[idx]) {
	case coverage.Partial:
		return "!"
	case coverage.Missed:
		return "!!"
	default:
		return ""
	}
}

// __IncByOne increments an integer by one
func __IncByOne(i int) int {
	return i + 1
}

// __GetRelativePath computes the relative path from one file to another
func __GetRelativePath(from, to string) string {
	depth := strings.Count(from, "/")
	if depth == 0 {
		return to
	}

	return strings.Repeat("../", depth) + to
}
