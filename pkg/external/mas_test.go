package external

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasImpl_List(t *testing.T) {
	mas := NewMas()

	output, err := mas.List(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
