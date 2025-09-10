package main

import "dagger/dag-coco/internal/dagger"

type ChangelogOptions struct {
	At         *string // changelog at point
	Template   *string // template
	Remote     *string // remote domain
	Owner      *string // repository owner
	Repository *string // repository
}

type CommitOptions struct {
	CogToml        *dagger.File // optional cog toml to overwrite file
	CommitType     string       // commit type: must be conv commit commit type
	CommitMessage  string       // commit message
	BreakingChange bool         // is a breaking change ?
}

type GetVersionOptions struct {
	Fallback *string // fallback version if not resolvable
	Package  *string // specific package
	Silence  bool    // only return version without additional info (-v)
}

type LogOptions struct {
	Authors      []string // filter authors
	CommitType   *string  // filter commit type
	Scope        *string  // filter scope
	NoError      bool     // Enforce (true) no-error
	BreakingOnly bool     // Show only breaking changes
}

type VersionBumpOptions struct {
	Auto      bool         // automatically choose next version
	CogToml   *dagger.File // provide a cog.toml config
	DryRun    bool         // dry-run
	Major     bool         // increment major
	Minor     bool         // increment minor
	Patch     bool         // increment patch
	PreMeta   *string      // Metadata used before building
	BuildMeta *string      // Metadata used when building
	Version   *string      // manually set version to (got prio compared to M.M.P SemVer)
	SkipCi    bool         // skip ci
}
