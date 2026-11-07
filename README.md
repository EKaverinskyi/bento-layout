# bento-layout

A standard folder layout for a portable drive, plus a small CLI to create and check it.

> **Status:** pre-alpha. The layout is still being designed and there is no code yet. See the [roadmap](ROADMAP.md).

## Idea

One folder on a drive holds all your portable data in a predictable, shell-friendly layout: lowercase names, no spaces, no special characters. The same layout on every drive and every machine.

A small CLI tool will help manage the layout: creating it, checking that a drive follows it, and, later, optionally integrating it with [zoxide](https://github.com/ajeetdsouza/zoxide) for quick navigation. Its exact commands and options will be defined in later phases of the [roadmap](ROADMAP.md).

## Layout

To be defined in `LAYOUT.md`. Designing it is the first phase of the [roadmap](ROADMAP.md).

## Non-goals

- Sync or backup engines (use rsync, restic or Syncthing)
- A UI
- Searching inside file contents or media
- Installing apps or system configuration
- Managing zoxide's own database

## Platforms

Primary: Linux on amd64 and macOS on Apple Silicon. Linux arm64 and macOS Intel are built on a best-effort basis.

## Documentation

- [PROJECT.md](PROJECT.md): project description and v0.1 contract
- [ROADMAP.md](ROADMAP.md): v1.0 definition and phased plan
- [CONTRIBUTING.md](CONTRIBUTING.md): how to work on the project
- [CHANGELOG.md](CHANGELOG.md): notable changes

## License

[MIT](LICENSE)
