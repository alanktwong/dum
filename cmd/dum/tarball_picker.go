package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type archivePicker func(context.Context, []string) (string, error)

func selectArchive(ctx context.Context, candidates []string) (string, error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("no archive candidates")
	}
	if len(candidates) == 1 {
		return candidates[0], nil
	}
	for _, candidate := range candidates {
		if strings.ContainsAny(candidate, "\r\n\x00") {
			return "", fmt.Errorf(
				"cannot interactively select an archive path containing a newline; pass the archive explicitly",
			)
		}
	}

	if fzfPath, err := exec.LookPath("fzf"); err == nil {
		return selectWithFZF(ctx, fzfPath, candidates)
	}
	if gumPath, err := exec.LookPath("gum"); err == nil {
		return selectWithGum(ctx, gumPath, candidates)
	}
	return "", fmt.Errorf("multiple tarballs, no fzf/gum:\n  %s", strings.Join(candidates, "\n  "))
}

func selectWithFZF(ctx context.Context, executable string, candidates []string) (string, error) {
	command := exec.CommandContext(ctx, executable, "--no-multi", "--prompt=restore> ")
	command.Stdin = strings.NewReader(strings.Join(candidates, "\n") + "\n")
	var output bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &output
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("fzf archive selection: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return validateArchiveSelection(output.String(), candidates)
}

func selectWithGum(ctx context.Context, executable string, candidates []string) (string, error) {
	arguments := append([]string{"choose"}, candidates...)
	command := exec.CommandContext(ctx, executable, arguments...)
	var output bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &output
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("gum archive selection: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return validateArchiveSelection(output.String(), candidates)
}

func validateArchiveSelection(selection string, candidates []string) (string, error) {
	selection = strings.TrimSuffix(selection, "\n")
	selection = strings.TrimSuffix(selection, "\r")
	if selection == "" {
		return "", fmt.Errorf("archive selection was empty")
	}
	for _, candidate := range candidates {
		if selection == candidate {
			return selection, nil
		}
	}
	return "", fmt.Errorf("picker returned an archive that was not a candidate: %q", selection)
}
