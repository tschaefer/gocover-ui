/*
Copyright (c) Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package coverage

import (
	"fmt"
	"math"
	"os"
	"strings"

	"golang.org/x/tools/cover"
)

// LineStatus represents the coverage status of a line
type LineStatus int

const (
	Missed LineStatus = iota
	Partial
	Covered
)

// String returns the string representation of LineStatus
func (ls LineStatus) String() (string, error) {
	switch ls {
	case Missed:
		return "missed", nil
	case Partial:
		return "partial", nil
	case Covered:
		return "covered", nil
	}

	return "", fmt.Errorf("invalid LineStatus: %d", ls)
}

// FileMetrics holds coverage metrics for a single file
type FileMetrics struct {
	FileName      string  `json:"fileName"`
	LocalPath     string  `json:"localPath"`
	TrackedLines  int     `json:"trackedLines"`
	CoveredLines  int     `json:"coveredLines"`
	PartialLines  int     `json:"partialLines"`
	MissedLines   int     `json:"missedLines"`
	CoveragePct   float64 `json:"coveragePct"`
	TotalStmts    int     `json:"totalStmts"`
	CoveredStmts  int     `json:"coveredStmts"`
	PerLineStatus []int   `json:"perLineStatus"`
}

type TotalMetrics struct {
	TotalFiles   int     `json:"totalFiles"`
	TotalStmts   int     `json:"totalStmts"`
	CoveredStmts int     `json:"coveredStmts"`
	CoveragePct  float64 `json:"coveragePct"`
}

// Analyze processes a coverage profile and returns file metrics
func Analyze(p *cover.Profile, module, srcRoot string) (*FileMetrics, error) {
	if !strings.HasPrefix(p.FileName, module) {
		return nil, fmt.Errorf("failed to match module %s", module)
	}

	localPath := strings.TrimPrefix(p.FileName, module+"/")

	source, err := os.ReadFile(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file %s: %w", localPath, err)
	}

	lines := strings.Split(string(source), "\n")
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	lineCount := len(lines)
	totalStmtsPerLine := make([]int, lineCount+1)
	coveredStmtsPerLine := make([]int, lineCount+1)

	totalStatements := 0
	coveredStatements := 0

	for _, b := range p.Blocks {
		for ln := b.StartLine; ln <= b.EndLine && ln <= lineCount; ln++ {
			totalStmtsPerLine[ln] += b.NumStmt
			if b.Count > 0 {
				coveredStmtsPerLine[ln] += b.NumStmt
			}
		}
		totalStatements += b.NumStmt
		if b.Count > 0 {
			coveredStatements += b.NumStmt
		}
	}

	trackedLines := 0
	coveredLines := 0
	partialLines := 0
	missedLines := 0
	perLineStatus := make([]int, lineCount+1)

	for ln := 1; ln <= lineCount; ln++ {
		total := totalStmtsPerLine[ln]
		if total == 0 {
			perLineStatus[ln] = -1
			continue
		}
		trackedLines++
		covered := coveredStmtsPerLine[ln]
		if covered == 0 {
			missedLines++
			perLineStatus[ln] = int(Missed)
		} else if covered >= total {
			coveredLines++
			perLineStatus[ln] = int(Covered)
		} else {
			partialLines++
			perLineStatus[ln] = int(Partial)
		}
	}

	coveragePct := 0.0
	if totalStatements > 0 {
		coveragePct = (float64(coveredStatements) / float64(totalStatements)) * 100.0
	}

	return &FileMetrics{
		FileName:      p.FileName,
		LocalPath:     localPath,
		TrackedLines:  trackedLines,
		CoveredLines:  coveredLines,
		PartialLines:  partialLines,
		MissedLines:   missedLines,
		CoveragePct:   round(coveragePct, 2),
		PerLineStatus: perLineStatus,
		TotalStmts:    totalStatements,
		CoveredStmts:  coveredStatements,
	}, nil
}

func Statistics(metrics []*FileMetrics) *TotalMetrics {
	totalFiles := len(metrics)
	totalStmts := 0
	coveredStmts := 0

	for _, m := range metrics {
		totalStmts += m.TotalStmts
		coveredStmts += m.CoveredStmts
	}

	coveragePct := 0.0
	if totalStmts > 0 {
		coveragePct = (float64(coveredStmts) / float64(totalStmts)) * 100.0
	}

	return &TotalMetrics{
		TotalFiles:   totalFiles,
		TotalStmts:   totalStmts,
		CoveredStmts: coveredStmts,
		CoveragePct:  round(coveragePct, 2),
	}
}

// Round rounds a float64 to specified precision
func round(v float64, prec int) float64 {
	p := math.Pow10(prec)
	return math.Round(v*p) / p
}
