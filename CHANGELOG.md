# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `check --layout <layout> --path <dir>` command and rule engine (ported from the Python prototype): both flags are required. Checks the layout against the universal rules (`--rules`, default `<layouts-dir>/rules.toml`) and its own rules and scenarios, then compares the real directory with it. Missing folders fail, extra folders are only listed. Exits 1 on any failure or incomplete layout.
- `init <layout> <target>` command: loads a layout (by path or by name from `--layouts-dir`) and creates its folders in the given target directory (required; use `.` for the current one); supports `--dry-run` and never changes existing folders.
- Layout loader (`internal/layout`): reads any layout directory with the same structure and validates its folder paths.
- Project description and high-level roadmap.
- Repository boilerplate: README, license, editor and git settings, markdown linting and pre-commit hooks.
- First layout definition, `personal`, in `layouts/personal/`: the folder tree, layout-specific rules and scenarios.
- Universal layout rules that every layout must follow, in `layouts/rules.toml`.
- Go project skeleton: a cobra CLI with a `version` command.
- Makefile with `build`, `test`, `lint`, `fmt` and `check` targets, and golangci-lint configuration with a pre-commit hook.
- GitHub Actions CI on every push: markdown checks, Go lint, workflow lint (actionlint) and tests on Linux and macOS. Binaries for all four target platforms are built on every push and kept as artifacts for 30 days.
- Releases: pushing a `v*.*.*` tag runs CI, then GoReleaser publishes archives and checksums for Linux and macOS (amd64 and arm64) as a GitHub Release.
