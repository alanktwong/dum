package external

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeImpl_InstallExtension(t *testing.T) {
	code := NewCode()

	err := code.InstallExtension(context.Background(), "esbenp.prettier-vscode")
	assert.NoError(t, err)
}

func TestCodeImpl_ListExtensions(t *testing.T) {
	code := NewCode()

	output, err := code.ListExtensions(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
