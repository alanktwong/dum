package plays

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayFactory_NoArgs(t *testing.T) {
	f := NewPlayFactory()
	assert.NotNil(t, f)
	assert.NotNil(t, f.Log)
	assert.NotNil(t, f.Utils)
	assert.NotNil(t, f.TaskFactory)
}


