package external

import (
	"fmt"
	"os"
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
	return absPath, nil
}
