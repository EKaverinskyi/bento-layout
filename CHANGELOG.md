# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Project description and high-level roadmap.
- Repository boilerplate: README, license, editor and git settings, markdown linting and pre-commit hooks.
- First layout definition, `personal`, in `layouts/personal/`: the folder tree, layout-specific rules and scenarios.
- Universal layout rules that every layout must follow, in `layouts/rules.toml`.
- Go project skeleton: a cobra CLI with a `version` command.
- Makefile with `build`, `test`, `lint`, `fmt` and `check` targets, and golangci-lint configuration with a pre-commit hook.
- GitHub Actions CI on every push: markdown checks, Go lint and tests on Linux and macOS. Binaries for all four target platforms are built on every push and, after CI passes, on version tags. Artifacts are kept 30 days.
