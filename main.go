package main

import (
	"context"
	"dagger/dag-coco/internal/dagger"
	"fmt"
)

type DagCoco struct{}

// Base Function returning a rust latest container with cocogitto installed
func (m *DagCoco) Base(ctx context.Context, repositoryUrl string, user string, gitToken *dagger.Secret) *dagger.Container {
	repository := m.CloneRepository(ctx, user, repositoryUrl, gitToken).Directory("/repository")

	return dag.Container().
		From("rust:latest").
		WithDirectory("/repository", repository).
		WithWorkdir("/repository").
		WithExec([]string{"cargo", "install", "--locked", "cocogitto"})
}

// Get the current version
// https://docs.cocogitto.io/guide/misc.html
func (m *DagCoco) GetVersion(ctx context.Context, repositoryUrl string, user string, gitToken *dagger.Secret, options *GetVersionOptions) (string, error) {

	args := ""
	if options != nil {
		if options.Package != nil {
			args = fmt.Sprintf("--package %s", *options.Package)
		}
		if options.Fallback != nil {
			args = fmt.Sprintf("%s --fallback %s", args, *options.Fallback)
		}
		if options.Silence {
			args = fmt.Sprintf("%s -v", args)
		}
	}

	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog get-version %s", args)}).
		Stdout(ctx)
}

// Create a new conventional commit using coco
// https://docs.cocogitto.io/guide/commit.html
func (m *DagCoco) Commit(ctx context.Context, repositoryUrl string, user string, gitToken *dagger.Secret, options *CommitOptions) (string, error) {

	base := m.Base(ctx, repositoryUrl, user, gitToken)
	args := ""

	if options == nil {
		return "", fmt.Errorf("Invalid commit options")
	}

	if options.CommitType == "" {
		return "", fmt.Errorf("Invalid commit type")
	}

	args = options.CommitType

	if options.BreakingChange {
		args = fmt.Sprintf("%s -B", args)
	}

	if options.CommitMessage == "" {
		return "", fmt.Errorf("Commit message can not be empty")
	}

	args = fmt.Sprintf("%s \"%s\"", args, options.CommitMessage)

	if options.CogToml != nil {
		base.WithFile("cog.toml", options.CogToml)
	}

	return base.Stdout(ctx)
}

// Check commit history against conventional commits spec
// https://docs.cocogitto.io/guide/check.html
func (m *DagCoco) Check(ctx context.Context, repositoryUrl string, user string, gitToken *dagger.Secret) (string, error) {

	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", "cog check"}).
		Stdout(ctx)
}

// Install a single or all hooks
// https://docs.cocogitto.io/guide/git_hooks.html
func (m *DagCoco) InstallHooks(ctx context.Context, repositoryUrl string, user string, gitToken *dagger.Secret, commits *string) (string, error) {

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

// Checks input string against conventional commits specification, error if not a conv commit
// https://docs.cocogitto.io/guide/verify.html
func (m *DagCoco) Verify(ctx context.Context, repositoryUrl string, user string, gitToken *dagger.Secret, commitMsg string) (string, error) {
	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog verify \"%s\"", commitMsg)}).
		Stdout(ctx)
}

// Get cocogitto log (like custom git log)
// https://docs.cocogitto.io/guide/log.html
func (m *DagCoco) Log(ctx context.Context, repositoryUrl string, user string, gitToken *dagger.Secret, version string, options *LogOptions) (string, error) {

	args := ""
	if options != nil {
		if options.Authors != nil {
			if len(options.Authors) != 0 {
				args = fmt.Sprintf("--author")
				for _, author := range options.Authors {
					args = fmt.Sprintf("%s \"%s\"", args, author)
				}
			}
		}

		if options.BreakingOnly {
			args = fmt.Sprintf("%s -B", args)
		}

		if options.CommitType != nil {
			args = fmt.Sprintf("%s --type %s", args, *options.CommitType)
		}

		if options.Scope != nil {
			args = fmt.Sprintf("%s --scope %s", args, *options.Scope)
		}

		if options.NoError {
			args = fmt.Sprintf("%s --no-error", args)
		}

	}

	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog log %s", args)}).
		Stdout(ctx)

}

// Generate a changelog automatically
// https://docs.cocogitto.io/guide/changelog.html
func (m *DagCoco) Changelog(ctx context.Context, repositoryUrl string, user string, gitToken *dagger.Secret, version string, options *ChangelogOptions) (string, error) {

	args := ""
	if options != nil {
		if options.At != nil {
			args = fmt.Sprintf("--at %s", *options.At)
		}

		if options.Template != nil {
			args = fmt.Sprintf("%s --template %s", args, *options.Remote)

			if *options.Template == "remote" {

				if options.Remote != nil {
					args = fmt.Sprintf("%s --remote %s", args, *options.Owner)
				}

				if options.Owner != nil {
					args = fmt.Sprintf("%s --owner %s", args, *options.Owner)
				}

				if options.Repository != nil {
					args = fmt.Sprintf("%s --repository %s", args, *options.Repository)
				}

			}
		}
	}
	return m.Base(ctx, repositoryUrl, user, gitToken).
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog changelog %s %s", args, version)}).
		Stdout(ctx)

}

// Bump version
// https://docs.cocogitto.io/guide/bump.html
func (m *DagCoco) Bump(ctx context.Context, repositoryUrl string, user string, gitToken *dagger.Secret, options *VersionBumpOptions) (string, error) {

	args := ""
	base := m.Base(ctx, repositoryUrl, user, gitToken)

	if options != nil {

		if options.Version == nil {
			if options.Auto {
				args = fmt.Sprintf("--auto")
			}

			if options.Major {
				args = fmt.Sprintf("%s --major", args)
			}

			if options.Minor {
				args = fmt.Sprintf("%s --minor", args)
			}

			if options.Patch {
				args = fmt.Sprintf("%s --patch", args)
			}
		} else {
			args = fmt.Sprintf("--version %s", *options.Version)
		}

		if options.CogToml != nil {
			base.WithFile("cog.toml", options.CogToml)
		}

		if options.DryRun {
			args = fmt.Sprintf("%s --dry-run", args)
		}

		if *options.PreMeta != "" {
			args = fmt.Sprintf("%s --pre \"%s\"", args, *options.PreMeta)
		}

		if *options.BuildMeta != "" {
			args = fmt.Sprintf("%s --build \"%s\"", args, *options.BuildMeta)
		}

		if options.SkipCi {
			args = fmt.Sprintf("%s --skip-ci", args)
		}

	}

	return base.
		WithExec([]string{"sh", "-c", fmt.Sprintf("cog bump %s", args)}).
		Stdout(ctx)
}
