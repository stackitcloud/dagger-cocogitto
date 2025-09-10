package main

import (
	"context"
	"dagger/dag-coco/internal/dagger"
	"fmt"
	"strings"
)

// Initializing Git
func (m *DagCoco) GitBase(ctx context.Context) *dagger.Container {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "upgrade"}).
		WithExec([]string{"apk", "add", "git"}).
		WithExec([]string{"git", "config", "--global", "user.name", "dagger"}).
		WithExec([]string{"git", "config", "--global", "user.email", "cicd@stackit.cloud"})
}

// Returns a container with teh cloned repository in /repository
func (m *DagCoco) CloneRepository(
	ctx context.Context,
	user string,
	repositoryUrl string,
	gitToken *dagger.Secret,
) *dagger.Container {
	token, err := gitToken.Plaintext(ctx)
	if err != nil {
		panic(err.Error())
	}

	url := fmt.Sprintf("https://%s:%s@%s", user, token, strings.TrimPrefix(repositoryUrl, "https://"))
	return m.GitBase(ctx).
		WithWorkdir("/repository").
		WithExec([]string{"git", "clone", url, "."})
}
