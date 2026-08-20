// Package rename provides file rename operations independent of Cobra.
package rename

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lg "alanktwong/dum/internal/logging"
)

// Options describes a single file rename operation.
type Options struct {
	Index       int
	FilePattern string
	Source      string
	Replace     string
	Lowercase   bool
	DryRun      bool
	TrimFront   int
	TrimBack    int
	Iterate     int
}

// Renamer performs file rename operations.
type Renamer struct {
	Log lg.Logger
}

// NewRenamer constructs a rename service.
func NewRenamer(logger lg.Logger) *Renamer {
	return &Renamer{Log: logger}
}

// Rename transforms and optionally renames one file.
func (s *Renamer) Rename(_ context.Context, options Options) error {
	const ellipsis = "..."
	pathname := filepath.Dir(options.FilePattern)
	s.Log.Debugf("%v pathname = %s", ellipsis, pathname)
	filename := strings.TrimSuffix(filepath.Base(options.FilePattern), filepath.Ext(options.FilePattern))
	s.Log.Debugf("%v filename = %s", ellipsis, filename)
	extension := filepath.Ext(options.FilePattern)
	s.Log.Debugf("%v extension = %s", ellipsis, extension)
	if options.TrimFront > 0 && options.TrimFront < len(filename) {
		filename = filename[options.TrimFront:]
		s.Log.Debugf("%v trimmed-front filename = %s", ellipsis, filename)
	}
	if options.TrimBack > 0 && options.TrimBack < len(filename) {
		filename = filename[:len(filename)-options.TrimBack]
		s.Log.Debugf("%v trimmed-back filename = %s", ellipsis, filename)
	}
	if options.Source != "" && options.Replace != "" {
		filename = strings.ReplaceAll(filename, options.Source, options.Replace)
		s.Log.Debugf("%v replaced filename = %s", ellipsis, filename)
	}
	if options.Lowercase {
		filename = strings.ToLower(filename)
		extension = strings.ToLower(extension)
		s.Log.Debugf("%v lowercased filename = %s", ellipsis, filename)
		s.Log.Debugf("%v lowercased extension = %s", ellipsis, extension)
	}
	if options.Iterate > 0 {
		filename = fmt.Sprintf("%s_%d", filename, options.Index+options.Iterate)
		s.Log.Debugf("%v iterated filename = %s", ellipsis, filename)
	}
	newPath := filepath.Join(pathname, filename+extension)
	s.Log.Infof("%v renaming %s -> %s", ellipsis, options.FilePattern, newPath)
	if !options.DryRun {
		if err := os.Rename(options.FilePattern, newPath); err != nil {
			return fmt.Errorf("error renaming file %s: %w", options.FilePattern, err)
		}
	}
	return nil
}
