// Package plays provides Play types and related utilities.
package plays

// PlayBookInfo provides information about a playbook.
type PlayBookInfo interface {
	GetID() string
	GetJetBrainsApps() map[string]string
}
