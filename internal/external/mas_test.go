package external

import (
	"context"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasImpl_List(t *testing.T) {
	if _, err := exec.LookPath("mas"); err != nil {
		t.Skip("mas is not installed")
	}

	mas := NewMas()

	output, err := mas.List(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
