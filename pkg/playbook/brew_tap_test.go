package playbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrewTap_Install(t *testing.T) {
	ctx := context.Background()
	tap := createTestBrewTap(t, Attributes{
		ID:          "homebrew/brew-tap",
		Description: "brew tap 1",
		Enabled:     true,
		Sudo:        false,
	})
	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().Tap(ctx, tap.ID).Return(nil).Once()
	tap.Brew = mockBrew
	// when
	input := createTestInput(t)
	input.DryRun = false
	got, err := tap.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, tap.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v): should install for happy path", input)
}

func TestBrewTap_DryRun(t *testing.T) {
	ctx := context.Background()
	tap := createTestBrewTap(t, Attributes{
		ID:          "homebrew/brew-tap",
		Description: "brew tap 1",
		Enabled:     true,
		Sudo:        false,
	})
	mockBrew := NewMockBrew(t)
	tap.Brew = mockBrew
	// when
	input := createTestInput(t)
	input.DryRun = true
	got, err := tap.Install(ctx, input)
	// then
	assert.Nil(t, err)
	wantedResult := expectTaskResult(t, true, tap.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v): should be unsuccessful for dry run", input)
	mockBrew.AssertNotCalled(t, "Tap")
}

func TestBrewTap_Disabled(t *testing.T) {
	ctx := context.Background()
	tap := createTestBrewTap(t, Attributes{
		ID:          "homebrew/brew-tap",
		Description: "brew tap 1",
		Enabled:     false,
		Sudo:        false,
	})
	mockBrew := NewMockBrew(t)
	tap.Brew = mockBrew
	// when
	input := createTestInput(t)
	input.DryRun = false
	got, err := tap.Install(ctx, input)
	// then
	assert.Nil(t, err)
	wantedResult := expectTaskResult(t, false, tap.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v): should be unsuccessful when disabled", input)
	mockBrew.AssertNotCalled(t, "Tap")
}

func TestBrewTap_Fail(t *testing.T) {
	ctx := context.Background()
	tap := createTestBrewTap(t, Attributes{
		ID:          "homebrew/brew-tap",
		Description: "brew tap 1",
		Enabled:     true,
		Sudo:        false,
	})
	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().Tap(ctx, tap.ID).Return(fmt.Errorf("fail")).Once()
	tap.Brew = mockBrew
	// when
	input := createTestInput(t)
	input.DryRun = false
	got, err := tap.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v): should fail install", input)
	assert.Nilf(t, got, "Install(%v): should fail install", input)
}
