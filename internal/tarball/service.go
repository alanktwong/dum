package tarball

import (
	"os"
	"time"

	lg "github.com/alanktwong/dum/internal/logging"
)

// Service owns tarball operations and their logging dependency.
type Service struct {
	Log     lg.Logger
	now     func() time.Time
	publish func(sourcePath, targetPath string) error
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

// ArtifactSet contains the archive and both required sidecars for one source.
type ArtifactSet struct {
	ArchivePath  string
	ChecksumPath string
	ManifestPath string
}

// BackupResult contains artifacts successfully created or planned by backup.
type BackupResult struct {
	Artifacts []ArtifactSet
}

// NewService constructs a tarball service.
func NewService(logger lg.Logger) *Service {
	return &Service{Log: logger, now: time.Now, publish: os.Link}
}
