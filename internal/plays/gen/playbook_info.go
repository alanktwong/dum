// Package gen provides generated types for plays package.
package gen

// PlayBookInfo provides information about a playbook.
type PlayBookInfo interface {
	GetID() string
	GetJetBrainsApps() map[string]string
}
