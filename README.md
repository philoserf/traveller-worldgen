# traveller-worldgen

A command-line generator for single Traveller worlds. Given a seed, it rolls a
world's Universal World Profile (UWP) — starport, size, atmosphere,
hydrographics, population, government, law level, tech level — plus its bases and
a generated name.

The repo is organized to hold **each edition of Traveller's world generation
side by side**. Today that is **Classic Traveller** (_Book 3: Worlds and
Adventures_, GDW 1977) in `classic/`, **MegaTraveller** (_Referee's Manual_,
DGP/GDW 1987) in `megatraveller/`, **Traveller: The New Era** (GDW 1993) in
`tne/`, and **Traveller5** (_Core Book 3_, FFE 5.10) in `t5/`; further editions
(Mongoose, …) are added as sibling packages sharing the same `dice` and `ehex`
foundation.

MegaTraveller extends the classic profile with a referee-chosen **subsector
nature** (which shapes the starport), a **non-imperial military base**, derived
**trade classifications**, and the two system-level counts on its generation
flowchart — **gas giants** and **planetoid belts**. TNE shares MegaTraveller's
nature-driven starport and military base but has no trade or system-count layer,
so a TNE world is a UWP plus bases.

Traveller5 is the outlier: its core draw is **Flux** (`1D − 1D`, −5…+5) rather
than the `2D − 7` of the other three, hydrographics keys off atmosphere instead
of size, and characteristics are hard-capped. On top of the UWP, T5 generates
Naval/Scout bases, the **PBG** counts (population digit, planetoid belts, gas
giants), trade classifications, and the **Ix/Ex/Cx** extensions (Importance,
Economic, Cultural).

## Usage

Each edition is a subcommand:

```
worldgen <edition> [flags]
worldgen classic [-seed N] [-format text|uwp|json] [-n COUNT]
worldgen mega    [-seed N] [-format text|uwp|json] [-n COUNT] [-nature NATURE]
worldgen tne     [-seed N] [-format text|uwp|json] [-n COUNT] [-nature NATURE]
worldgen t5      [-seed N] [-format text|uwp|json] [-n COUNT]
```

Run `worldgen` with no arguments to list editions. Flags shared by every
edition:

- `-seed N` — seed for reproducibility (omitted → time-based).
- `-format` — `text` (default), `uwp`, or `json`.
- `-n COUNT` — generate COUNT independent worlds (default 1). A single seed
  reproduces the whole batch.

MegaTraveller and TNE add:

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
Pupijou  X200346-2  None
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
  Belts         1

$ go run ./cmd/worldgen mega -seed 1977 -n 4 -format uwp
Paedrucil      C645856-7  SM  G3  -   -
Niagaevei      B452452-B  N   G3  -   Ni Po
Voupaher       C748754-9  -   -   -   Ag
Daehia         E787645-5  -   G4  B1  Ag Ni Ri

$ go run ./cmd/worldgen tne -seed 1977 -n 4 -format uwp
Robriheil      C645856-7  SM
Sihir          B675453-A  —
Tifeifarn      C496333-5  S
Liafufous      A83A336-D  N

$ go run ./cmd/worldgen t5 -seed 42
Saejiarn  X223300-3  None
  Starport      X  None; no starport
  Size          2  2,000 miles (3,200 km)
  Atmosphere    2  Very Thin, Tainted
  Hydrographics 3  30% Water
  Population    3  Thousands (10^3)
  Government    0  No Government Structure
  Law Level     0  No Law. No prohibitions
  Tech Level    3  Pre-Industrial (basic science)
  Trade         Lo Po
  Importance    {-3}
  Economic      (620+1)  (RU 12)
  Cultural      [2161]
  Gas Giants    2
  Belts         3

$ go run ./cmd/worldgen t5 -seed 1977 -n 4 -format uwp
Cirisould      C623468-A  -  {+0} (A33+3) [648C]  130  Ni Px Po
Faedoth        B486899-7  S  {+1} (978+2) [A914]  821  Ph Pa Ri
Mumeil         B221244-D  S  {+1} (B11-3) [434D]  201  Lo Po
Niadreicoun    B898200-7  -  {-1} (310+1) [5188]  701  Lo
```

The `uwp` base column is a compact code. Classic: `N` naval, `S` scout, `NS`
both, `—` none. MegaTraveller and TNE: `N` naval, `S` scout, `A` naval + scout,
plus a trailing `M` for a non-imperial military base. The MegaTraveller line then
carries a gas-giant marker (`G`_n_), a planetoid-belt marker (`B`_n_), and the
trade codes (`-` for an empty field); TNE's line ends at the base code (`—` when
none). Traveller5 uses `N`/`S`/`A` bases (`-` when none), then the extensions
`{Ix}` `(Ex)` `[Cx]`, the three-digit **PBG**, and the trade codes.

## Layout

Shared foundation (edition-independent):

- `dice/` — the `Roller` interface (`Seeded`, `Scripted`, `Fixed`) and `D6(n)`.
- `ehex/` — Traveller extended-hex digit `Encode`/`Decode`.

Per-edition rules (each self-contained, sharing `dice` and `ehex`):

- `classic/` — Classic Traveller Book 3: tables, the tech-index DM matrix, world
  generation, descriptions, and the name generator.
- `megatraveller/` — MegaTraveller Referee's Manual: the same shape plus the
  nature-driven starport table, bases (incl. military), trade classifications,
  gas giants, and planetoid belts.
- `tne/` — Traveller: The New Era: MegaTraveller's UWP generation (nature-driven
  starport, military base) minus the trade and gas-giant layers.
- `t5/` — Traveller5 Core Book 3: the Flux-based UWP, Naval/Scout bases, PBG
  counts, trade classifications, and the Ix/Ex/Cx extensions.
- _future_ — `mongoose/`, … as sibling packages.

CLI and docs:

- `cmd/worldgen/` — the CLI. `main.go` dispatches on an edition subcommand;
  each edition has its own runner file (`classic.go`, `megatraveller.go`,
  `tne.go`, `t5.go`) registered in the `editions` map.
- `docs/<edition>/` — the source PDF(s) and the verified rules extract.

**Adding an edition** means: a new rules package (e.g. `mongoose/`), a
`cmd/worldgen/<edition>.go` runner, and one entry in the `editions` map. The
rules packages stay deliberately independent (duplication over premature
abstraction), but the CLI scaffolding is shared: a generic `runEdition[T]` owns
the flag parse / seed / generate / write flow, and `renderWorlds[T]` /
`renderWorldsJSON[T]` own output formatting — so a runner supplies only its
edition-specific flags, generator, and per-world formatters. A cross-edition
`Generator` interface over the rules packages is still **not** defined; it hasn't
earned one.

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
`B` is deferred — the generator emits `N`/`S`/`A`/`M`. Gas giants and planetoid
belts (the flowchart's two system-level counts) are generated; the fuller
Extended System Generation (stars, orbits, satellites, additional planets),
travel zones, and allegiance remain out of scope.

### TNE (The New Era)

**Source.** The _World Tamer's Handbook_ (which prompted this edition) turned out
not to contain the core UWP generation — its Stars & Planets chapter defers to
"pages 192–195 of the TNE rulebook." So the tables come from the **TNE core
rulebook** (Basic Mainworld Generation, printed pp. 186–190); the WTH's
"Detailing Planets" physical layer is out of scope.

**Near-identical to MegaTraveller.** TNE's UWP generation matches MegaTraveller's
almost exactly — same starport-by-nature table, formulas, and bases. It notably
prints **Starport A tech-DM `+6`**, confirming the value MegaTraveller misprints
as `+8`. The only rules differences: the tech-DM matrix adds Starport F `+1` and
Government E/F `−1`, and TNE has **no trade-classification or gas-giant step**.

**Deferred.** TNE's Planetary-Density → surface-gravity step (part of world size)
and the WTH detail layer (temperature, atmospheric composition, world mass) are
physical data beyond the UWP, out of scope.

### Traveller5 (Core Book 3)

**A different engine.** T5's core draw is **Flux** (`1D − 1D`, −5…+5), not
`2D − 7`. Atmosphere is `Flux + Size`, hydrographics is `Flux + Atmosphere`
(keyed off atmosphere, not size), government/law are Flux-based, and every
characteristic is hard-capped (Size/Atm/Gov `0–F`, Law `0–J`, Hyd `0–A`) — so
nothing is ever "beyond described range." The starport table is a single 2D
column identical to Classic; bases roll `2D ≤ target` (T5's `N−` notation), the
mirror of Classic's `≥`.

**Extensions.** Beyond the UWP, T5 generates the **Ix** (Importance,
`{±n}`), **Ex** (Economic, `(RLI±E)` with derived Resource Units), and **Cx**
(Cultural, `[HASS]`) extensions, plus the **PBG** counts and trade
classifications.

**Source note.** Core Book 3 p. 25 gives the tech-level DM formula but no prose
TL table, so the era-band descriptions reuse the edition-invariant Traveller TL
bands. Book 3 states no As→Va suppression rule, so an asteroid belt carries both.

**Deferred.** MainWorld-type (planet vs. satellite), Habitable-Zone/Climate and
the climate trade codes that depend on it, the Secondary/Political/Special trade
codes (referee-assigned), Nobility/Allegiance/Travel-Zone/Native-Intelligent-Life
(the rest of NABZ-NIL), Depot/Way-Station bases, and the whole system map are out
of scope for a bare mainworld generator.
