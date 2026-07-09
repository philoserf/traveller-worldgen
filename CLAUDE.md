# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
task              # default: gofumpt check + go vet + golangci-lint + go test -race
task check        # lint/format/vet only (no tests)
task test         # go test -race ./...
task fmt          # apply gofumpt -extra (formatting is CI-enforced)
task lint:fix     # golangci-lint --fix

go test ./megatraveller -run TestWorkedExample   # a single test
go run ./cmd/worldgen mega -seed 42               # run the CLI
```

Formatting uses **gofumpt with `-extra`** (not plain `gofmt`); the `fmt:check`
task fails CI on any diff. Lint is golangci-lint v2 (config in `.golangci.yml`,
notably `prealloc`, `revive`, `unparam`, `gocritic`).

**Golden files** live in `cmd/worldgen/testdata/`. After an intentional change
to CLI output, regenerate them: `go test ./cmd/worldgen -run <Edition> -update`
(e.g. `-run Mega`), then eyeball the diff before committing.

The `gopls` MCP server is registered in `.mcp.json` (project scope). Use its
tools (`go_workspace`, `go_search`, `go_file_context`, `go_diagnostics`) for
navigation and diagnostics once approved.

## Architecture

A CLI that rolls single Traveller worlds. The organizing principle is **one Go
package per Traveller edition, sitting side by side on a shared foundation** —
so understanding any single edition means understanding this shared shape.

- `dice/` and `ehex/` are **edition-independent** and reused verbatim by every
  edition. `dice.Roller` is the seam through which all randomness flows;
  `ehex` is Traveller extended-hex (0-9, A-Z minus I/O).
- `classic/` and `megatraveller/` are **self-contained rules packages** — each
  has its own `world.go` (the `World` struct, `UWP()`, `*Desc()` methods,
  `MarshalJSON`), `generate.go` (`Generate`), `tables.go` (all lookup maps +
  tech-DM matrices), and `names.go`. Duplication between editions is
  **deliberate**: they evolve independently, so `names.go` and the small
  `floor0`/`clamp`/generic-`lookup` helpers are copied per package rather than
  shared. Do not "DRY up" the rules packages into a common one.
- `cmd/worldgen/` is the CLI. `main.go` dispatches on the first arg via the
  `editions` map; each edition has one runner file (`classic.go`,
  `megatraveller.go`). **Rendering lives in the CLI, not the packages** — the
  `text`/`uwp` layout, column widths, and sentinels are presentation concerns.
  The only genuinely shared CLI helper is the generic `renderWorldsJSON[T]`.

**Adding an edition**: a new rules package, a `cmd/worldgen/<edition>.go` runner,
and one entry in the `editions` map. That's the whole contract.

### The determinism contract (critical)

`Generate(r dice.Roller, …)` draws dice in a **fixed, documented order**, and
**consumes no die for a branch that doesn't apply** (bases on E/X starports, the
automatic size-0/1 results, the gas-giant quantity when none are present). This
is what lets a seed reproduce an entire `-n` batch, and it's load-bearing for the
tests. If you reorder draws or add/remove a conditional die, golden files and the
scripted-dice oracle tests will shift — that's expected; regenerate and re-derive.

`dice.Scripted` (used by oracle tests) returns one preset value **per `D6` call
regardless of `n`** — the value is the pre-computed 2D sum. So a test sequence is
read as "the natural roll for each successive draw," in the exact draw order.

### Rules provenance and the "keep values raw" philosophy

Every table is transcribed from `docs/<edition>/world-generation.md`, which is
verified cell-by-cell against the scanned source PDFs (symlinked beside it).
**That markdown is the source of truth for table values** — change it there, then
mirror into `tables.go`. Where the printed source contradicts itself, the doc
records the resolution (e.g. MegaTraveller's Technology-DM table prints Starport A
as `+8`, but the manual's stated max roll of 20 only reconciles at `+6`, so `+6`
is used).

Characteristics are **floored at 0 but otherwise kept raw** and rendered in
extended hex, so a DM-inflated value can exceed its last table row. Classic
reports such values as `"(beyond described range)"`; MegaTraveller's ranges are
complete (Gov/Atmo 0-F, Law/Tech 0-L) so it never needs that. The lone real cap
is **hydrographics at A** (100% water). Law level is special: its description
always appends the enforcement saving throw (`N+ to avoid arrest`), and its raw
(uncapped) value is derived from the true uncapped government so the two UWP
digits stay consistent.

### Test styles

Four complementary styles per the classic/mega suites: a **scripted-dice oracle**
(`TestWorkedExample` — hand-computed sequence → exact `World`), **known-value
table tests** (`*Desc()`, trade codes), **property sweeps** (hundreds of seeds
asserting ranges/determinism/well-formed UWP), and **golden CLI tests**
(`cmd/worldgen`, `-update`-refreshable). When adding an edition, mirror all four.
