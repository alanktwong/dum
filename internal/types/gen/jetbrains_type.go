package gen

//go:generate go run github.com/abice/go-enum@latest --file=$GOFILE --names

// JetBrainsType is an enum source.
// ENUM(clion, datagrip, goland, idea, phpstorm, pycharm, rider, rubymine, rustrover, webstorm).
type JetBrainsType string
