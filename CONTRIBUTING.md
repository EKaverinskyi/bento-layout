# Contributing

## Setup

1. Install [lefthook](https://lefthook.dev/installation/) (`brew install lefthook`) and Node.js (for markdown checks).
2. Install the latest stable [Go](https://go.dev/dl/); `cmd/go.mod` sets the minimum version.
3. Install [golangci-lint](https://golangci-lint.run/) v2 built with that Go version or newer: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`. A binary built with an older Go (such as an outdated Homebrew bottle) refuses to run.
4. Run `lefthook install` once in the repository.

Commits that fail the markdown or Go checks are blocked locally.

## Go code

The Go module lives in `cmd/`; run `go` commands from there. From the repository root, use `make`:

| Command      | What it does                                    |
| ------------ | ----------------------------------------------- |
| `make build` | Build the `bento-layout` binary                 |
| `make test`  | Run tests with the race detector                |
| `make lint`  | Run golangci-lint (config: `cmd/.golangci.yml`) |
| `make fmt`   | Format Go code                                  |
| `make check` | Run lint and tests                              |

Run `make help` for the full list. The pre-commit hook runs only the linter; run `make test` yourself before opening a pull request.

## Checks

Run the same checks manually:

```sh
npx --yes prettier@3.6.2 --check "**/*.md"
npx --yes markdownlint-cli2@0.18.1
```

Fix formatting automatically with `npx --yes prettier@3.6.2 --write "**/*.md"`.

## Conventions

- Branch from `main`; open a pull request for every change.
- Keep commits focused. Use short imperative subjects, for example `Add manifest verification`.
- Update [CHANGELOG.md](CHANGELOG.md) under `[Unreleased]` for user-visible changes.
- Folder and file names in this repository are lowercase, without spaces or special characters.
- The public [ROADMAP.md](ROADMAP.md) lists goals only.
