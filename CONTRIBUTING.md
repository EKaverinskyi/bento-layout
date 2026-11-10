# Contributing

## Setup

1. Install [lefthook](https://lefthook.dev/installation/) (`brew install lefthook`) and Node.js (for markdown checks).
2. Install the latest stable [Go](https://go.dev/dl/); `go.mod` sets the minimum version.
3. Install [golangci-lint](https://golangci-lint.run/) v2 (`brew install golangci-lint`). It refuses to run if it was built with an older Go minor version than the one in `go.mod`; in that case upgrade it or build it yourself: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`.
4. Run `lefthook install` once in the repository.

Commits that fail the markdown or Go checks are blocked locally.

## Go code

The Go module is at the repository root: `cmd/bento-layout` holds `main` and `internal/` holds the packages. Use `make` or plain `go` commands:

| Command      | What it does                                |
| ------------ | ------------------------------------------- |
| `make build` | Build the `bento-layout` binary             |
| `make test`  | Run tests with the race detector            |
| `make lint`  | Run golangci-lint (config: `.golangci.yml`) |
| `make fmt`   | Format Go code                              |
| `make check` | Run lint and tests                          |

Run `make help` for the full list. The pre-commit hook runs only the linter; run `make test` yourself before opening a pull request.

## Checks

Run the same checks manually:

```sh
npx --yes prettier@3.6.2 --check "**/*.md"
npx --yes markdownlint-cli2@0.18.1
```

Fix formatting automatically with `npx --yes prettier@3.6.2 --write "**/*.md"`.

## Continuous integration

GitHub Actions runs on every push, to any branch:

- `CI` (`ci.yml`): markdown format and lint, workflow lint (actionlint), `go mod tidy` check and golangci-lint, then `go test -race` on Linux and macOS. Tests run only after the checks pass.
- `Build` (`build.yml`): builds binaries for Linux and macOS (amd64 and arm64) and keeps them as workflow artifacts for 30 days, so you can download and try a branch.

## Conventions

- Branch from `main`; open a pull request for every change.
- Keep commits focused. Use short imperative subjects, for example `Add manifest verification`.
- Update [CHANGELOG.md](CHANGELOG.md) under `[Unreleased]` for user-visible changes.
- Folder and file names in this repository are lowercase, without spaces or special characters.
- The public [ROADMAP.md](ROADMAP.md) lists goals only.
