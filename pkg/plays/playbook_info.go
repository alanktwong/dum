package plays

type PlayBookInfo interface {
	GetID() string
	GetJetBrainsApps() map[string]string
}
