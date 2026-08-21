package external

import (
	"context"
	"os/exec"
	"testing"

	"github.com/alanktwong/dum/internal/external/gen"
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

func TestMasImpl_List_UnsupportedOS(t *testing.T) {
	ext := gen.NewMockExt(t)
	ext.On("IsOSX").Return(false)
	mas := &MasImpl{ext: ext}

	output, err := mas.List(context.Background())
	assert.Error(t, err)
	assert.Empty(t, output)
	assert.Contains(t, err.Error(), "outside of macOS")
}

func TestMasImpl_Install_UnsupportedOS(t *testing.T) {
	ext := gen.NewMockExt(t)
	ext.On("IsOSX").Return(false)
	mas := &MasImpl{ext: ext}

	err := mas.Install(context.Background(), "123456789")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "outside of macOS")
}
