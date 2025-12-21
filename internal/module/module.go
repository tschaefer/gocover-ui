/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package module

import (
	"fmt"
	"os"
	"path/filepath"
)

// Read parses the go.mod file located at path and returns the module name.
func Read(path string) (string, error) {
	goModFile, err := os.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod file: %w", err)
	}

	var module string
	_, err = fmt.Sscanf(string(goModFile), "module %s", &module)
	if err != nil {
		return "", fmt.Errorf("failed to parse module name: %w", err)
	}

	return module, nil
}
