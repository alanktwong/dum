package installer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithSudo(t *testing.T) {
	ctx := context.Background()
	ctx = WithSudo(ctx, true)
	got := GetSudo(ctx)
	assert.True(t, got)
}

func TestWithSudo_False(t *testing.T) {
	ctx := context.Background()
	ctx = WithSudo(ctx, false)
	got := GetSudo(ctx)
	assert.False(t, got)
}

func TestGetSudo_NotSet(t *testing.T) {
	ctx := context.Background()
	got := GetSudo(ctx)
	assert.False(t, got)
}

func TestGetSudo_WrongType(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, sudoCtxKey{}, "not a bool")
	got := GetSudo(ctx)
	assert.False(t, got)
}
