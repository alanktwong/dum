package tarball

import lg "github.com/alanktwong/dum/internal/logging"

// Service owns tarball operations and their logging dependency.
type Service struct {
	Log lg.Logger
}

// BackupOptions describes one backup invocation after CLI input resolution.
type BackupOptions struct {
	Paths  []string
	Target string
	Home   string
	DryRun bool
}

// RestoreOptions describes one archive restore operation.
type RestoreOptions struct {
	Archive          string
	Target           string
	Destination      string
	Overwrite        bool
	AllowUnsafeLinks bool
	PreserveOwner    bool
	DryRun           bool
}

// NewService constructs a tarball service.
func NewService(logger lg.Logger) *Service {
	return &Service{Log: logger}
}
