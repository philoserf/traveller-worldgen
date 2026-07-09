# traveller-worldgen

A command-line generator for single Traveller worlds. Given a seed, it rolls a
world's Universal World Profile (UWP) — starport, size, atmosphere,
hydrographics, population, government, law level, tech level — plus its bases and
a generated name.

The repo is organized to hold **each edition of Traveller's world generation
side by side**. Today that is **Classic Traveller** (_Book 3: Worlds and
Adventures_, GDW 1977) in the `classic/` package; future editions
(MegaTraveller, Mongoose, T5, …) are added as sibling packages sharing the same
`dice` and `ehex` foundation.

## Usage

Each edition is a subcommand:

```
worldgen <edition> [flags]
worldgen classic [-seed N] [-format text|uwp|json] [-n COUNT]
```

Run `worldgen` with no arguments to list editions. Classic flags:

- `-seed N` — seed for reproducibility (omitted → time-based).
- `-format` — `text` (default), `uwp`, or `json`.
- `-n COUNT` — generate COUNT independent worlds (default 1). A single seed
  reproduces the whole batch.

### Examples

```
$ go run ./cmd/worldgen classic -seed 42
Draejapu  X200346-2  None
  Starport      X  No starport; no provision for starship landings
  Size          2  2,000 miles diameter
  Atmosphere    0  No atmosphere
  Hydrographics 0  No free standing water
  Population    3  1,000
  Government    4  Representative Democracy
  Law Level     6  Most firearms (all except shotguns) prohibited; 6+ to avoid arrest
  Tech Level    2  Halberd, Broadsword · Cannon · Wind

$ go run ./cmd/worldgen classic -seed 1977 -n 4 -format uwp
Horou          C674474-9  S
Peimern        E7A5431-5  —
Nojeis         C110875-A  S
Houreir        C250555-9  —

$ go run ./cmd/worldgen classic -seed 42 -format json
{
  "name": "Draejapu",
  "uwp": "X200346-2",
  ...
}
```

The `uwp` base column is a compact code: `N` naval, `S` scout, `NS` both, `—`
none.

## Layout

Shared foundation (edition-independent):

- `dice/` — the `Roller` interface (`Seeded`, `Scripted`, `Fixed`) and `D6(n)`.
- `ehex/` — Traveller extended-hex digit `Encode`/`Decode`.

Per-edition rules:

- `classic/` — Classic Traveller Book 3: tables, the tech-index DM matrix, world
  generation, descriptions, and the name generator.
- _future_ — `mega/`, `mongoose/`, `t5/`, … each self-contained, sharing `dice`
  and `ehex`.

CLI and docs:

- `cmd/worldgen/` — the CLI. `main.go` dispatches on an edition subcommand;
  each edition has its own runner file (`classic.go`) registered in the
  `editions` map.
- `docs/classic/` — the Book 3 source PDF and the verified rules extract.

**Adding an edition** means: a new rules package (e.g. `mega/`), a
`cmd/worldgen/mega.go` runner, and one entry in the `editions` map. A
cross-edition `Generator` interface and a shared renderer are deliberately
**not** defined yet — the right seam only becomes clear with a second edition in
hand, so they will be extracted then rather than guessed at now.

## Development

```
task          # gofumpt + go vet + golangci-lint + go test -race
task test     # tests only
task fmt      # apply gofumpt
```

Golden CLI outputs live in `cmd/worldgen/testdata/`; regenerate them with
`go test ./cmd/worldgen -update` after an intentional output change.

## Rules & provenance

Rules and tables are transcribed from `docs/classic/world-generation.md`, an extraction
verified cell-by-cell against the Book 3 source PDF.

**Clamping.** Book 3 floors negative dice-modifier results at 0 but specifies no
upper bounds. This generator floors every characteristic at 0 and otherwise
keeps the rolled value, rendering it in extended hex — so atmosphere,
government, law level, and tech level may all exceed their last described table
row. The sole upper cap is **hydrographics at A (100% water)**, a real physical
ceiling. Keeping the other values raw matches the source: the Book 3 tech-index
matrix itself defines DMs for atmosphere D and E, and law level is deliberately
derived from the true (uncapped) government so both digits stay consistent.

**Describing out-of-range values.** Atmosphere, government, and tech level have
no Book 3 guidance past their tables, so a value beyond the last row is reported
as "(beyond described range)". **Law level is the exception**: Book 3's note
makes each level cumulative and ties the raw level to an enforcement throw, so
the text output always shows the weapons prohibition plus the saving throw to
avoid arrest (`N+`, equal to the law level). A law level above 9 keeps level
9's prohibition while the throw keeps rising — becoming impossible to make (and
thus certain arrest) at 13+.
