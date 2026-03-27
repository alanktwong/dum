package external

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (u *DefaultExt) absolutePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("file path is empty ")
	}
	absPath, err := u.ExpandUser(path)
	if err != nil {
		return "", err
	}
	if strings.Contains(absPath, "$") {
		absPath = os.ExpandEnv(absPath)
	}
	if !filepath.IsAbs(absPath) {
		absPath, err = filepath.Abs(absPath)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path for %s: %w", path, err)
		}
	}
	return absPath, nil
}
