package external

import (
	"alanktwong/dum/internal/external/gen"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	joined := filepath.Join(os.Getenv("HOME"),
		"Library",
		"Application Support",
		"JetBrains",
		ideName,
		"plugins",
		plugin)
	pluginPath := filepath.Clean(joined)
	if _, err := os.Stat(pluginPath); err != nil {
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
