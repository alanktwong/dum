package playbook

import "context"

// Lister can list given an input.
type Lister interface {
	List(ctx context.Context, input *Input) (*TaskResult, error)
}

// Installer can install given an input.
type Installer interface {
	Install(ctx context.Context, input *Input) (*TaskResult, error)
}

// Task can list and install given an input.
type Task interface {
	Lister
	Installer
}
