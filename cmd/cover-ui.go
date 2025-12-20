package coverui

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	ver         = flag.Bool("version", false, "print version and exit")
)

func Run() {
	flag.Parse()

	if *ver {
		version.Print()
		return
	}

	module, err := module.Read(*srcRoot)
	checkErr(err)

	profiles, err := cover.ParseProfiles(*profileFile)
	checkErr(err)

	var files []*coverage.FileMetrics
	for _, p := range profiles {
		fm, err := coverage.Analyze(p, module, *srcRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping %s: %v\n", p.FileName, err)
			continue
		}
		files = append(files, fm)
	}

	if len(files) == 0 {
		checkErr(fmt.Errorf("no valid coverage profiles found in %s", *profileFile))
	}

	if *cleanOutDir {
		err := os.RemoveAll(*outDir)
		checkErr(err)
	}

	filesDir := filepath.Join(*outDir, "tree")
	err = file.Assets(filesDir)
	checkErr(err)

	for _, f := range files {
		err := file.Generate(f, filesDir)
		checkErr(err)
	}

	err = index.Assets(*outDir)
	checkErr(err)

	err = index.Generate(files, *outDir, module)
	checkErr(err)
}

func checkErr(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
