package main

import (
	"context"
	"dagger/dag-coco/internal/dagger"
	"fmt"
	"strings"
)

type DagCoco struct{}

// Base Function returning a rust latest container with cocogitto installed
func (m *DagCoco) Base(
	// current context
	ctx context.Context,
	// remote repository url
	repositoryUrl string,
	// user account name
	user string,
	// user account git token
	gitToken *dagger.Secret) *dagger.Container {
	repository := m.CloneRepository(ctx, user, repositoryUrl, gitToken).Directory("/repository")

	return dag.Container().
		From("rust:latest").
		WithDirectory("/repository", repository).
		WithWorkdir("/repository").
		WithExec([]string{"cargo", "install", "--locked", "cocogitto"})
}

// Bump version
// https://docs.cocogitto.io/guide/bump.html
func (m *DagCoco) Bump(
	// context
	ctx context.Context,
	// remote repository URL
	repositoryUrl string,
	// remote username
	user string,
	// secret Git token
	gitToken *dagger.Secret,
	// +optional
	// auto bump version
	auto bool,
	// +optional
	// cogToml file reference
	cogToml *dagger.File,
	// +optional
	// should run dry
	dryRun bool,
	// +optional
	// increment major
	major bool,
	// +optional
	// increment minor
	minor bool,
	// +optional
	// increment patch
	patch bool,
	// +optional
	// metadata before building
	preMeta string,
	// +optional
	// build metadata
	buildMeta string,
	// +optional
	// set specfic version (prio over major,minor,patch)
	version string,
	// +optional
	// Skip ci flag
	skipCi bool,
) (string, error) {

	args := ""
	base := m.Base(ctx, repositoryUrl, user, gitToken)

	if version == "" {
		if auto {
			args = "--auto"
		}

		if major {
			args = fmt.Sprintf("%s --major", args)
		}

		if minor {
			args = fmt.Sprintf("%s --minor", args)
		}

		if patch {
			args = fmt.Sprintf("%s --patch", args)
		}
	} else {
		args = fmt.Sprintf("--version %s", version)
	}

	if cogToml != nil {
		base.WithFile("cog.toml", cogToml)
	}

	if dryRun {
		args = fmt.Sprintf("%s --dry-run", args)
	}

	if preMeta != "" {
		args = fmt.Sprintf("%s --pre \"%s\"", args, preMeta)
	}

	if buildMeta != "" {
		args = fmt.Sprintf("%s --build \"%s\"", args, buildMeta)
	}

	if skipCi {
		args = fmt.Sprintf("%s --skip-ci", args)
	}

	result, err := base.
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog bump %s", args)}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("error bumping version")
	}

	return strings.TrimSuffix(result, "\n"), err

}

// Generate a changelog automatically
// https://docs.cocogitto.io/guide/changelog.html
func (m *DagCoco) Changelog(
	// context vontext
	ctx context.Context,
	// remote repo address
	repositoryUrl string,
	// remote user account
	user string,
	// secret git token
	gitToken *dagger.Secret, version string,
	// +optional
	// changelog between 1.00 and at
	at string,
	// +optional
	// changelog template: default, full_hash, remote
	template string,
	// +optional
	// remote domain: e.g. google.com
	remote string,
	// +optional
	// user account / organization
	owner string,
	// +optional
	// repository name
	repository string,
) (string, error) {

	args := ""

	if at != "" {
		args = fmt.Sprintf("--at %s", at)
	}

	if template != "" {
		args = fmt.Sprintf("%s --template %s", args, template)

		if template == "remote" {

			if remote != "" {
				args = fmt.Sprintf("%s --remote %s", args, remote)
			}

			if owner != "" {
				args = fmt.Sprintf("%s --owner %s", args, owner)
			}

			if repository != "" {
				args = fmt.Sprintf("%s --repository %s", args, repository)
			}

		}
	}

	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog changelog %s %s", args, version)}).
		Stdout(ctx)

}

// Check commit history against conventional commits spec
// https://docs.cocogitto.io/guide/check.html
func (m *DagCoco) Check(
	// current context
	ctx context.Context,
	// remote repository url
	repositoryUrl string,
	// user account
	user string,
	// remote git token secret
	gitToken *dagger.Secret,
	// +optional
	cogToml *dagger.File) (string, error) {

	base := m.Base(ctx, repositoryUrl, user, gitToken)

	if cogToml != nil {
		base = base.WithFile("cog.toml", cogToml)
	}

	return base.
		WithExec([]string{"sh", "-c", "cog check"}).
		Stdout(ctx)
}

// Create a new conventional commit using coco
// https://docs.cocogitto.io/guide/commit.html
func (m *DagCoco) Commit(
	// current context
	ctx context.Context,
	// remote repo url
	repositoryUrl string,
	// user account
	user string,
	// repository secret token
	gitToken *dagger.Secret,
	// +optional
	// Cocogitto config file
	cogToml *dagger.File,
	// +optional
	// tyoe of commit
	commitType string,
	// +optional
	// commit message
	commitMessage string,
	// +optional
	// is it a breaking change
	breakingChange bool,
) (string, error) {

	base := m.Base(ctx, repositoryUrl, user, gitToken)
	args := ""

	if commitType == "" {
		return "", fmt.Errorf("invalid commit type")
	}

	args = commitType

	if breakingChange {
		args = fmt.Sprintf("%s -B", args)
	}

	if commitMessage == "" {
		return "", fmt.Errorf("commit message can not be empty")
	}

	args = fmt.Sprintf("%s \"%s\"", args, commitMessage)

	if cogToml != nil {
		base.WithFile("cog.toml", cogToml)
	}

	return base.WithExec([]string{"sh", "-c", fmt.Sprintf("cog commit %s", args)}).Stdout(ctx)
}

// Get the current version
// https://docs.cocogitto.io/guide/misc.html
func (m *DagCoco) GetVersion(
	// current context
	ctx context.Context,
	// remote repo URL
	repositoryUrl string,
	// username
	user string,
	// secret git token
	gitToken *dagger.Secret,
	// +optional
	// fallback version
	fallback string,
	// +optional
	// define package name
	pack string,
	// +optional
	// Only get version
	silence bool,
) (string, error) {

	args := ""

	if pack != "" {
		args = fmt.Sprintf("--package %s", pack)
	}
	if fallback != "" {
		args = fmt.Sprintf("%s --fallback %s", args, fallback)
	}
	if silence {
		args = fmt.Sprintf("-v %s", args)
	}

	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog get-version %s", args)}).
		Stdout(ctx)
}

// Install a single or all hooks
// https://docs.cocogitto.io/guide/git_hooks.html
func (m *DagCoco) InstallHooks(
	// the current context
	ctx context.Context,
	// remote repo URL
	repositoryUrl string,
	// username
	user string,
	// git secret token
	gitToken *dagger.Secret,
	// semver, commit, range commit xxx..xxy
	commits *string) (string, error) {

	args := ""
	if commits != nil {
		args = *commits
	} else {
		args = "--all"
	}

	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog install-hook %s", args)}).
		Stdout(ctx)
}

// Get cocogitto log (like custom git log)
// https://docs.cocogitto.io/guide/log.html
func (m *DagCoco) Log(
	// current context
	ctx context.Context,
	// repository remote address
	repositoryUrl string,
	// username
	user string,
	// git token
	gitToken *dagger.Secret,
	// +optional
	// filter authors
	authors []string,
	// +optional
	// define commit type: e.g. feat
	commitType string,
	// +optional
	// scope of commit
	scope string,
	// +optional
	// no error enforce
	noError bool,
	// +optional
	// filter breaking changes only
	breakingOnly bool,
) (string, error) {

	args := ""

	if authors != nil {
		if len(authors) != 0 {
			args = "--author"
			for _, author := range authors {
				args = fmt.Sprintf("%s \"%s\"", args, author)
			}
		}
	}

	if breakingOnly {
		args = fmt.Sprintf("%s -B", args)
	}

	if commitType != "" {
		args = fmt.Sprintf("%s --type %s", args, commitType)
	}

	if scope != "" {
		args = fmt.Sprintf("%s --scope %s", args, scope)
	}

	if noError {
		args = fmt.Sprintf("%s --no-error", args)
	}

	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog log %s", args)}).
		Stdout(ctx)

}

// Checks input string against conventional commits specification, error if not a conv commit
// https://docs.cocogitto.io/guide/verify.html
func (m *DagCoco) Verify(
	// current context
	ctx context.Context,
	// repository url
	repositoryUrl string,
	// username
	user string,
	// git token secret
	gitToken *dagger.Secret,
	// message to check
	commitMsg string) (string, error) {
	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog verify \"%s\"", commitMsg)}).
		Stdout(ctx)
}
