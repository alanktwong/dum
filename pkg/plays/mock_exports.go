package plays

import (
	"testing"

	pg "awong/dotfiles/pkg/plays/gen"
)

type MockPlayBookInfo = pg.MockPlayBookInfo

func NewMockPlayBookInfo(t *testing.T) *pg.MockPlayBookInfo {
	return pg.NewMockPlayBookInfo(t)
}
