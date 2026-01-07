/*
Copyright (c) Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package tree

import (
	"math"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tschaefer/cover-ui/internal/coverage"
)

// Node represents a directory or file in the file browser
type Node struct {
	Name         string                `json:"name"`
	Path         string                `json:"path"`
	IsDir        bool                  `json:"isDir"`
	File         *coverage.FileMetrics `json:"file,omitempty"`
	Children     []*Node               `json:"children,omitempty"`
	TrackedLines int                   `json:"trackedLines"`
	CoveredLines int                   `json:"coveredLines"`
	PartialLines int                   `json:"partialLines"`
	MissedLines  int                   `json:"missedLines"`
	TotalStmts   int                   `json:"totalStmts"`
	CoveredStmts int                   `json:"coveredStmts"`
	CoveragePct  float64               `json:"coveragePct"`
}

// Build creates a hierarchical tree structure from flat file list
func Build(files []*coverage.FileMetrics) *Node {
	root := &Node{
		Name:     "/",
		Path:     "/",
		IsDir:    true,
		Children: make([]*Node, 0),
	}

	for _, file := range files {
		current := root
		currentPath := ""

		pathParts := strings.Split(file.LocalPath, "/")
		for i := 0; i < len(pathParts)-1; i++ {
			part := pathParts[i]
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = filepath.Join(currentPath, part)
			}

			found := false
			for _, child := range current.Children {
				if child.Name == part && child.IsDir {
					current = child
					found = true
					break
				}
			}

			if found {
				continue
			}

			dirNode := createDirNode(part, currentPath)
			current.Children = append(current.Children, dirNode)
			current = dirNode
		}

		fileNode := createFileNode(file, pathParts[len(pathParts)-1])
		current.Children = append(current.Children, fileNode)
	}

	calculateDirCoverage(root)

	sortNodes(root)

	return root
}

// createDirNode creates a directory node
func createDirNode(name, path string) *Node {
	return &Node{
		Name:     name,
		Path:     path,
		IsDir:    true,
		Children: make([]*Node, 0),
	}
}

// createFileNode creates a file node
func createFileNode(fileMetrics *coverage.FileMetrics, name string) *Node {
	return &Node{
		Name:         name,
		Path:         fileMetrics.FileName,
		IsDir:        false,
		File:         fileMetrics,
		TrackedLines: fileMetrics.TrackedLines,
		CoveredLines: fileMetrics.CoveredLines,
		PartialLines: fileMetrics.PartialLines,
		MissedLines:  fileMetrics.MissedLines,
		TotalStmts:   fileMetrics.TotalStmts,
		CoveredStmts: fileMetrics.CoveredStmts,
		CoveragePct:  fileMetrics.CoveragePct,
	}
}

// calculateDirCoverage recursively calculates coverage stats for directories
func calculateDirCoverage(node *Node) {
	if !node.IsDir {
		return
	}

	totalTracked := 0
	totalCovered := 0
	totalPartial := 0
	totalMissed := 0
	totalStmts := 0
	coveredStmts := 0

	for _, child := range node.Children {
		if child.IsDir {
			calculateDirCoverage(child)
		}
		totalTracked += child.TrackedLines
		totalCovered += child.CoveredLines
		totalPartial += child.PartialLines
		totalMissed += child.MissedLines
		totalStmts += child.TotalStmts
		coveredStmts += child.CoveredStmts
	}

	node.TrackedLines = totalTracked
	node.CoveredLines = totalCovered
	node.PartialLines = totalPartial
	node.MissedLines = totalMissed
	node.TotalStmts = totalStmts
	node.CoveredStmts = coveredStmts
	if totalStmts > 0 {
		node.CoveragePct = round((float64(coveredStmts)/float64(totalStmts))*100.0, 2)
	}
}

// Round rounds a float to specified precision
func round(v float64, prec int) float64 {
	p := math.Pow10(prec)
	return math.Round(v*p) / p
}

// sortNodes sorts tree nodes: directories first, then alphabetically
func sortNodes(node *Node) {
	if !node.IsDir {
		return
	}

	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		return node.Children[i].Name < node.Children[j].Name
	})

	for _, child := range node.Children {
		if child.IsDir {
			sortNodes(child)
		}
	}
}
