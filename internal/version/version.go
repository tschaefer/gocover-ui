/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package version

import (
	"fmt"
	"os"
)

var (
	GitCommit, Version string
)

func Release() string {
	if Version == "" {
		Version = "dev"
	}

	return Version
}

func Commit() string {
	return GitCommit
}

func Banner() string {
	return `
  __ _  ___   ___ _____   _____ _ __      _   _(_)
 / _' |/ _ \ / __/ _ \ \ / / _ \ '__|____| | | | |
| (_| | (_) | (_| (_) \ V /  __/ | |_____| |_| | |
 \__, |\___/ \___\___/ \_/ \___|_|        \__,_|_|
 |___/
`
}

func Print() {
	no_color := os.Getenv("NO_COLOR")
	if no_color != "" {
		fmt.Printf("%s\n", Banner())
	} else {
		fmt.Printf("\033[34m%s\033[0m\n", Banner())
	}
	fmt.Printf("Release: %s\n", Release())
	fmt.Printf("Commit:  %s\n", Commit())
}
