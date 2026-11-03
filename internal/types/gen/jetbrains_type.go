package gen

//go:generate go-enum --file=$GOFILE --names

// JetBrainsType is an enum source.
// ENUM(clion, datagrip, goland, idea, phpstorm, pycharm, rider, rubymine, rustrover, webstorm).
type JetBrainsType string
