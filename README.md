# traveller-worldgen

A command-line generator for single Traveller worlds. Given a seed, it rolls a
world's Universal World Profile (UWP) — starport, size, atmosphere,
hydrographics, population, government, law level, tech level — plus its bases and
a generated name.

The repo is organized to hold **each edition of Traveller's world generation
side by side**. Today that is **Classic Traveller** (_Book 3: Worlds and
Adventures_, GDW 1977) in the `classic/` package and **MegaTraveller**
(_Referee's Manual_, DGP/GDW 1987) in the `megatraveller/` package; further
editions (Mongoose, T5, …) are added as sibling packages sharing the same `dice`
and `ehex` foundation.

MegaTraveller extends the classic profile with a referee-chosen **subsector
nature** (which shapes the starport), a **non-imperial military base**, derived
**trade classifications**, and a **gas-giant** count.

## Usage

Each edition is a subcommand:

```
worldgen <edition> [flags]
worldgen classic [-seed N] [-format text|uwp|json] [-n COUNT]
worldgen mega    [-seed N] [-format text|uwp|json] [-n COUNT] [-nature NATURE]
```

Run `worldgen` with no arguments to list editions. Flags shared by every
edition:

- `-seed N` — seed for reproducibility (omitted → time-based).
- `-format` — `text` (default), `uwp`, or `json`.
- `-n COUNT` — generate COUNT independent worlds (default 1). A single seed
  reproduces the whole batch.

MegaTraveller adds:

- `-nature` — subsector nature selecting the starport column: `backwater`,
  `standard` (default), `mature`, or `cluster`.

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

$ go run ./cmd/worldgen mega -seed 42
Japu  X200346-2  None
  Starport      X  None; no starport
  Size          2  Small (3,200 km, Luna)
  Atmosphere    0  Vacuum
  Hydrographics 0  Desert World (0–4%)
  Population    3  Low (thousands)
  Government    4  Representative Democracy
  Law Level     6  Moderate law (all firearms except shotguns prohibited); 6+ to avoid arrest
  Tech Level    2  Pre-Industrial (printing press)
  Trade         Lo Ni Va
  Gas Giants    5

$ go run ./cmd/worldgen mega -seed 1977 -n 4 -format uwp
Briheil        C645856-7  SM  G3  -
Hite           B675453-A  -   G1  Ni
Tifeifarn      E328755-7  -   -   -
Fufougaess     A83A336-D  N   G3  Lo Ni Wa
```

The `uwp` base column is a compact code. Classic: `N` naval, `S` scout, `NS`
both, `—` none. MegaTraveller: `N` naval, `S` scout, `A` naval + scout, plus a
trailing `M` for a non-imperial military base (`-` when empty); its line then
carries a gas-giant marker (`G`_n_) and the trade codes, with `-` for an empty
field.

## Layout

Shared foundation (edition-independent):

- `dice/` — the `Roller` interface (`Seeded`, `Scripted`, `Fixed`) and `D6(n)`.
- `ehex/` — Traveller extended-hex digit `Encode`/`Decode`.

Per-edition rules (each self-contained, sharing `dice` and `ehex`):

- `classic/` — Classic Traveller Book 3: tables, the tech-index DM matrix, world
  generation, descriptions, and the name generator.
- `megatraveller/` — MegaTraveller Referee's Manual: the same shape plus the
  nature-driven starport table, bases (incl. military), trade classifications,
  and gas giants.
- _future_ — `mongoose/`, `t5/`, … as sibling packages.

CLI and docs:

- `cmd/worldgen/` — the CLI. `main.go` dispatches on an edition subcommand;
  each edition has its own runner file (`classic.go`, `megatraveller.go`)
  registered in the `editions` map.
- `docs/<edition>/` — the source PDF(s) and the verified rules extract.

**Adding an edition** means: a new rules package (e.g. `mongoose/`), a
`cmd/worldgen/<edition>.go` runner, and one entry in the `editions` map. The
rules packages stay deliberately independent (duplication over premature
abstraction); the only shared CLI seam extracted so far is the generic
`renderWorldsJSON[T]`. A cross-edition `Generator` interface is still **not**
defined — rendering differs enough per edition that it hasn't earned one.

## Development

```
task          # gofumpt + go vet + golangci-lint + go test -race
task test     # tests only
task fmt      # apply gofumpt
```

Golden CLI outputs live in `cmd/worldgen/testdata/`; regenerate them with
`go test ./cmd/worldgen -update` after an intentional output change.

## Rules & provenance

Each edition's rules and tables are transcribed from
`docs/<edition>/world-generation.md`, an extraction verified cell-by-cell
against the source PDF(s).

### Classic (Book 3)

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

### MegaTraveller (Referee's Manual)

**Complete ranges.** Unlike Book 3, MegaTraveller describes every value the dice
can produce (government/atmosphere `0–F`, law/tech `0–L`), so nothing is ever
"beyond described range." Hydrographics is still capped at A.

**Source corrections.** Two printed-table errors were resolved against the
manual's own text and recorded in the doc: the Technology-DM table prints
Starport A as `+8`, but its stated maximum adjusted roll of 20 reconciles only at
`+6` (the value used); and the trade-classification atmosphere ranges for
Industrial (`2–4,7,9`) and Rich (`6,8`) were corrected from an earlier
reconstruction.

**Deferred.** The Basic Mainworld Generation flowchart has no bases step and the
sources give no throw that promotes a scout base to a Way Station, so base code
`B` is deferred — the generator emits `N`/`S`/`A`/`M`. Extended System
Generation (orbits, satellites, planetoid belts, travel zones, allegiance) is out
of scope.
