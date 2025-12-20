package tree

import (
	"testing"

	"github.com/tschaefer/cover-ui/internal/coverage"
)

func TestBuild(t *testing.T) {
	files := []*coverage.FileMetrics{
		{
			FileName:     "github.com/test/repo/main.go",
			LocalPath:    "main.go",
			TrackedLines: 10,
			CoveredLines: 8,
			PartialLines: 1,
			MissedLines:  1,
			TotalStmts:   10,
			CoveredStmts: 8,
			CoveragePct:  80.0,
		},
		{
			FileName:     "github.com/test/repo/pkg/utils.go",
			LocalPath:    "pkg/utils.go",
			TrackedLines: 20,
			CoveredLines: 15,
			PartialLines: 2,
			MissedLines:  3,
			TotalStmts:   20,
			CoveredStmts: 15,
			CoveragePct:  75.0,
		},
		{
			FileName:     "github.com/test/repo/internal/cmd/cmd.go",
			LocalPath:    "internal/cmd/cmd.go",
			TrackedLines: 15,
			CoveredLines: 10,
			PartialLines: 2,
			MissedLines:  3,
			TotalStmts:   15,
			CoveredStmts: 10,
			CoveragePct:  66.67,
		},
	}

	root := Build(files)

	if root == nil {
		t.Fatal("Expected non-nil root node")
	}
	if !root.IsDir {
		t.Error("Expected root to be a directory")
	}
	if len(root.Children) == 0 {
		t.Fatal("Expected root to have children")
	}

	if root.TrackedLines != 45 {
		t.Errorf("Expected root TrackedLines 45, got %d", root.TrackedLines)
	}
	if root.CoveredLines != 33 {
		t.Errorf("Expected root CoveredLines 33, got %d", root.CoveredLines)
	}
}

func TestNodeSorting(t *testing.T) {
	files := []*coverage.FileMetrics{
		{
			FileName:     "github.com/test/repo/zebra.go",
			LocalPath:    "zebra.go",
			TrackedLines: 5,
			TotalStmts:   5,
			CoveredStmts: 5,
		},
		{
			FileName:     "github.com/test/repo/pkg/utils.go",
			LocalPath:    "pkg/utils.go",
			TrackedLines: 10,
			TotalStmts:   10,
			CoveredStmts: 10,
		},
		{
			FileName:     "github.com/test/repo/alpha.go",
			LocalPath:    "alpha.go",
			TrackedLines: 3,
			TotalStmts:   3,
			CoveredStmts: 3,
		},
	}

	root := Build(files)

	if len(root.Children) == 5 {
		t.Fatalf("Expected at 5 children, got %d", len(root.Children))
	}

	if !root.Children[0].IsDir {
		t.Error("Expected first child to be a directory")
	}

	fileNames := []string{}
	for _, child := range root.Children {
		if !child.IsDir {
			fileNames = append(fileNames, child.Name)
		}
	}

	if len(fileNames) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(fileNames))
	}

	if fileNames[0] > fileNames[1] {
		t.Errorf("Expected files to be sorted alphabetically, got %v", fileNames)
	}
}
