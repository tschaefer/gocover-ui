/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package coverui

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/tools/cover"

	"github.com/tschaefer/cover-ui/internal/coverage"
	"github.com/tschaefer/cover-ui/internal/generator/file"
	"github.com/tschaefer/cover-ui/internal/generator/index"
	"github.com/tschaefer/cover-ui/internal/module"
	"github.com/tschaefer/cover-ui/internal/version"
)

var (
	profileFile = flag.String("profile", "coverage.out", "coverage profile file")
	outDir      = flag.String("out", "coverage", "output directory for generated html files")
	srcRoot     = flag.String("src", ".", "source root directory on disk")
	cleanOutDir = flag.Bool("clean", false, "clean output directory before generating files")
	versionInfo = flag.Bool("version", false, "print version and exit")
	quiet       = flag.Bool("quiet", false, "suppress progress and statistics output")
)

func Run() {
	flag.Parse()

	printVersion()

	removeOldFiles()

	module, err := module.Read(*srcRoot)
	checkErr(err)

	files := analyzeProfile(module)
	generateHtmlFiles(files, module)
}

func printVersion() {
	if !*versionInfo {
		return
	}

	version.Print()
	os.Exit(0)
}

func removeOldFiles() {
	if !*cleanOutDir {
		return
	}

	err := os.RemoveAll(*outDir)
	checkErr(err)
}

func analyzeProfile(module string) []*coverage.FileMetrics {
	profiles, err := cover.ParseProfiles(*profileFile)
	checkErr(err)

	var files []*coverage.FileMetrics
	for _, profile := range profiles {
		metrics, err := coverage.Analyze(profile, module, *srcRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping %s: %v\n", profile.FileName, err)
			continue
		}
		files = append(files, metrics)
	}

	if len(files) == 0 {
		checkErr(fmt.Errorf("no valid coverage profiles found in %s", *profileFile))
	}

	return files
}

func generateHtmlFiles(files []*coverage.FileMetrics, module string) {
	filesDir := filepath.Join(*outDir, "tree")

	err := file.Assets(filesDir)
	checkErr(err)

	err = index.Assets(*outDir)
	checkErr(err)

	err = index.Generate(files, *outDir, module)
	checkErr(err)

	for _, f := range files {
		err := file.Generate(f, filesDir)
		checkErr(err)

		path := filepath.Join(filesDir, strings.TrimSuffix(f.LocalPath, ".go")+".html")
		printProgress(fmt.Sprintf("Generated %s", path))
	}

	printStatistics(files)
}

func printProgress(message string) {
	if *quiet {
		return
	}

	fmt.Printf("%s", message)
	time.Sleep(125 * time.Millisecond)
	fmt.Print("\x1b[2K\r")
}

func printStatistics(files []*coverage.FileMetrics) {
	if *quiet {
		return
	}

	totalMetrics := coverage.Statistics(files)
	fmt.Printf(
		"Generated coverage report for %.2f%% of statements (%d of %d) in %d files.\n",
		totalMetrics.CoveragePct,
		totalMetrics.CoveredStmts,
		totalMetrics.TotalStmts,
		totalMetrics.TotalFiles,
	)
}

func checkErr(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
