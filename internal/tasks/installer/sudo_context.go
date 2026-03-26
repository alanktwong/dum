package installer

import "context"

// sudoCtxKey is a private type key for the sudo value that may be in context.Context.
type sudoCtxKey struct{}

// WithSudo adds sudo to context.
func WithSudo(ctx context.Context, sudo bool) context.Context {
	return context.WithValue(ctx, sudoCtxKey{}, sudo)
}

// GetSudo gets sudo from context.
func GetSudo(ctx context.Context) bool {
	sudo, ok := ctx.Value(sudoCtxKey{}).(bool)
	if !ok {
		return false
	}
	return sudo
}
