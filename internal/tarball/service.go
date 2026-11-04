package tarball

import (
	"os"
	"time"

	lg "github.com/alanktwong/dum/internal/logging"
)

// Service owns tarball operations and their logging dependency.
type Service struct {
	Log         lg.Logger
	now         func() time.Time
	publish     func(sourcePath, targetPath string) error
	changeOwner func(root *os.Root, path string, uid, gid int, symlink bool) error
	ownerOf     func(os.FileInfo) (uid, gid int64, ok bool)
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

// RestoreResult describes the archive entries validated or restored.
type RestoreResult struct {
	ArchivePath string
	Destination string
	Paths       []string
}

// NewService constructs a tarball service.
func NewService(logger lg.Logger) *Service {
	return &Service{
		Log:         logger,
		now:         time.Now,
		publish:     os.Link,
		changeOwner: changeOwnerDefault,
		ownerOf:     fileOwner,
	}
}
