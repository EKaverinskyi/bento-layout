# Roadmap

## Vision

A single, memorable, shell-friendly folder convention for a portable drive, plus a small, dependable command-line tool that helps create and maintain it.

## Principles

- **The layout is the product.** The tool serves the layout, not the other way around.
- **Shell-friendly by design.** Names work in a terminal without quoting or escaping.
- **Covers real life.** The layout has a home for every kind of file the user actually keeps.
- **Navigation helpers come later.** Optional conveniences, such as zoxide integration, are built after the core tool works.
- **Small and boring.** One static binary, no services, no UI.

## Platforms

Primary: Linux on amd64 and macOS on Apple Silicon. Other targets are built on a best-effort basis.

## Phases

| Phase                 | Goal                                                                          | Outcome                                                                                                                                                                                             |
| --------------------- | ----------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Layout design         | Decide the standard layout.                                                   | A written specification that can be recited from memory and that the user is willing to move their files onto. Frozen before any layout code is written.                                            |
| Foundations           | A project where every change is automatically checked.                        | Formatting, linting and tests run automatically on every change, locally and in CI.                                                                                                                 |
| Core tool (v0.1)      | Create a layout, check a drive against it, and notice changes to stored data. | A layout can be created safely on an empty or populated drive, and checked with clear, actionable results. Changes to stored data are never accepted silently. Binaries and a README are published. |
| Extras and polish     | Optional conveniences, and a tool that is pleasant for daily use.             | Helpers such as zoxide integration work, and nothing else depends on them. Good help, shell completions and documentation.                                                                          |
| Real-world use        | Run the layout and tool on the user's real drive.                             | Daily use for a meaningful period, with findings recorded.                                                                                                                                          |
| Stable release (v1.0) | Fix what real use revealed and commit to stability.                           | Layout, configuration format and behaviour are stable and documented.                                                                                                                               |

Layout design and Foundations are independent and can overlap. Everything after them depends on the frozen layout.

## Definition of done for v1.0

- The layout is frozen, documented, and in daily use on the author's real drive without further renaming of top-level folders.
- The core tool is tested and works on both primary platforms.
- Releases are built and published automatically.
- Documentation explains the layout, the reasoning behind it, and the project's non-goals.

## Non-goals

See the [README](README.md#non-goals).
