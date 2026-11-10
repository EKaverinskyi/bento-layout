# bento-layout

A standard folder layout for a portable drive, plus a small CLI to create and check it.

**Not to be confused with "bento grid" UI design.** `bento-layout` is a tool for folder layout management. The name comes from the bento box: one container with fixed compartments.

> **Status:** pre-alpha. The layout is still being designed. The CLI is only a skeleton (it can print its version) and layout creation and checking are not implemented yet. See the [roadmap](ROADMAP.md).

## Idea

One folder on a drive holds all your portable data in a predictable, shell-friendly layout: lowercase names, no spaces, no special characters. The same layout on every drive and every machine.

A small CLI tool (written in Go) will help manage the layout: creating it, checking that a drive follows it, and, later, optionally integrating it with [zoxide](https://github.com/ajeetdsouza/zoxide) for quick navigation. Its exact commands and options will be defined in later phases of the [roadmap](ROADMAP.md).

## Layout

The layout is described as data: a layout is a folder of TOML files in [`layouts/`](layouts/). The first one, `personal`, is a draft and is not frozen yet. The frozen specification will be written in `LAYOUT.md`; designing it is the first phase of the [roadmap](ROADMAP.md).

How layouts work:

- Every layout is self-contained: its folders, its own rules and a set of scenarios (an item and the one folder it belongs in).
- Universal rules in [`layouts/rules.toml`](layouts/rules.toml) apply to every layout. For example, folder names are lowercase, and sibling folders can be told apart by tab completion within four characters.
- You use one layout at a time. Moving to another layout is up to you.
- You can extend your layout over time, for example by adding `projects/3d-modelling`.
- The tool will compare a drive with the layout and warn about differences. It will not change anything without your confirmation. It only looks at names, extensions and metadata, never at file content.

Top-level folders of the `personal` draft:

- `dev`: programming projects you develop, split by purpose
- `projects`: non-programming projects, such as 3D or video
- `programs`: software you use but do not develop
- `games`: game installs, saves and settings
- `media`: standalone pictures, music, audio and videos
- `books`: books, PDFs and papers
- `documents`: contracts, notes and scans
- `config`: dotfiles and application configuration
- `downloads`: files downloaded from the internet
- `rebuildable`: re-downloadable bulk data
- `unsorted`: your own inbox

## Non-goals

- Sync or backup engines (use rsync, restic or Syncthing)
- A UI
- Searching inside file contents or media
- Installing apps or system configuration
- Managing zoxide's own database

## Platforms

Primary: Linux on amd64 and macOS on Apple Silicon. Linux arm64 and macOS Intel are built on a best-effort basis.

## Documentation

- [ROADMAP.md](ROADMAP.md): v1.0 definition and phased plan
- [CONTRIBUTING.md](CONTRIBUTING.md): how to work on the project
- [CHANGELOG.md](CHANGELOG.md): notable changes

## License

[MIT](LICENSE)
