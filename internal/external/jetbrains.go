package external

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alanktwong/dum/internal/external/gen"
)

// JetBrainsApp is an abstraction for jetbrains app.
type JetBrainsApp = gen.JetBrainsApp

// JetBrainsImpl implements JetBrainsApp.
type JetBrainsImpl struct{}

// NewJetBrains constructs JetBrainsImpl.
func NewJetBrains() *JetBrainsImpl {
	return &JetBrainsImpl{}
}

// IsInstalled implements JetBrainsApp.
func (j *JetBrainsImpl) IsInstalled(ideName, plugin string) bool {
	base := filepath.Join(os.Getenv("HOME"),
		"Library",
		"Application Support",
		"JetBrains")
	joined := filepath.Join(base,
		ideName,
		"plugins",
		plugin)
	pluginPath := filepath.Clean(joined)
	rel, err := filepath.Rel(base, pluginPath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false
	}
	if _, err := os.Stat(pluginPath); err != nil { // #nosec G703 -- rel check above confines pluginPath to base.
		// ctx.Log.Debugf("%v %v: plugin %v already installed for %s", TASK_ELLIPSIS, ctx.Play, app, t.ID)
		return true
	}
	return false
}

// Install implements JetBrainsApp.
func (j *JetBrainsImpl) Install(ctx context.Context, app, plugin string) error {
	cmd := exec.CommandContext(ctx, app, "installPlugins", plugin)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to %s installPlugins %s: %w", app, plugin, err)
	}
	return nil
}
