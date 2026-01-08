/*
Copyright (c) Tobias Schäfer. All rights reserved.
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

	err := removeOldFiles()
	checkErr(err)

	module, err := module.Read(*srcRoot)
	checkErr(err)

	files, err := analyzeProfile(module)
	checkErr(err)

	err = generateHtmlFiles(files, module)
	checkErr(err)
}

func printVersion() {
	if !*versionInfo {
		return
	}

	version.Print()
	os.Exit(0)
}

func removeOldFiles() error {
	if !*cleanOutDir {
		return nil
	}

	return os.RemoveAll(*outDir)
}

func analyzeProfile(module string) ([]*coverage.FileMetrics, error) {
	profiles, err := cover.ParseProfiles(*profileFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coverage profile: %w", err)
	}

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
		return nil, fmt.Errorf("no valid coverage profiles found in %s", *profileFile)
	}

	return files, nil
}

func generateHtmlFiles(files []*coverage.FileMetrics, module string) error {
	filesDir := filepath.Join(*outDir, "tree")

	if err := file.Assets(filesDir); err != nil {
		return err
	}

	if err := index.Assets(*outDir); err != nil {
		return err
	}

	if err := index.Generate(files, *outDir, module); err != nil {
		return err
	}

	for _, f := range files {
		if err := file.Generate(f, filesDir); err != nil {
			return err
		}

		path := filepath.Join(filesDir, strings.TrimSuffix(f.LocalPath, ".go")+".html")
		printProgress(fmt.Sprintf("Generated %s", path))
	}

	printStatistics(files)

	return nil
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
