# Contributing

## Setup

1. Install [lefthook](https://lefthook.dev/installation/) (`brew install lefthook`) and Node.js (for markdown checks).
2. Run `lefthook install` once in the repository.

Commits that fail the markdown format or lint checks are blocked locally.

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
