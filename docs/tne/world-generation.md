# Traveller: The New Era — Basic Mainworld Generation

Extracted from _Traveller: The New Era_ (Game Designers' Workshop, 1993), the
"Imperial Universal World Profile (UWP) Generation" flowchart (printed p. 186)
and the "Universal World Profile Tables" / "World Building Charts" (printed
pp. 188–190).

Scope: the **mainworld** only (single world), matching the `classic` and
`megatraveller` packages. Extended System Generation (stars, orbits, satellites,
gas giants, planetoid belts) is out of scope, as is the _World Tamer's Handbook_
"Detailing Planets" physical layer (see § Out of scope).

> **Source note.** The _World Tamer's Handbook_ (the supplement that prompted
> this edition) does **not** contain the core UWP generation — its "Stars &
> Planets" chapter (p. 84) states that planets are generated "according to the
> system found on pages 192–195 of the TNE rulebook," which "provides the
> standard UWP data." So the tables below come from the **TNE core rulebook**
> (printed pp. 186–190). The WTH adds only a detail layer on top (temperature,
> atmospheric composition/pressure/taint, world mass, mapping), captured here as
> an out-of-scope extension.

---

## How TNE differs from MegaTraveller

TNE's Basic Mainworld Generation is **nearly identical to MegaTraveller's** — the
same starport-by-nature table, the same core formulas, the same bases table, and
the same law/tech description tables. The differences are small:

- **Technology DM matrix** adds three cells MegaTraveller lacks: Starport **F =
  +1**, Government **E = −1**, Government **F = −1**. Starport A is printed as
  **+6** (unlike the MegaTraveller Referee's Manual, which misprints it as +8).
- **Atmosphere** codes D/E/F are labelled as Exotic variants (Exotic dense-high /
  Exotic thin-low / Exotic ellipsoid).
- **World Size** step also rolls a Planetary Density (1D6) used to derive surface
  gravity — physical detail beyond the UWP, deferred here.
- **Bases** and combined codes (N/S/A/B/M) are identical to MegaTraveller,
  including the deferred Way Station (code `B`, which has no separate throw).

Because the systems are so close, the `tne` package mirrors `megatraveller`
closely; the divergences above are the only rules changes.

---

## Generation sequence (mainworld)

Steps 1–2 (subsector mapping, system details) are out of scope for a single
mainworld. The mainworld characteristics:

| Step | Characteristic | Throw                                                               |
| ---- | -------------- | ------------------------------------------------------------------- |
| 3    | Starport       | 2D → [Starport table](#starport-2d) for the chosen subsector nature |
| 4    | World Size     | 2D − 2 (a 1D density → surface-gravity step is deferred)            |
| 5    | Atmosphere     | 2D − 7 + Size (Size 0 → Atmosphere 0)                               |
| 6    | Hydrographics  | 2D − 7 + Size (Size ≤ 1 → 0; Atmosphere ≤ 1 or ≥ A → DM −4)         |
| 7    | Population     | 2D − 2                                                              |
| 8    | Government     | 2D − 7 + Population                                                 |
| 9    | Law Level      | 2D − 7 + Government                                                 |
| 10   | Tech Level     | 1D + [tech DMs](#technology-dm-matrix), floored at 0                |
| 11   | Bases          | per-starport throws → N/S/A/B/M                                     |

Adjusted rolls below 0 are treated as 0. As in MegaTraveller, Government/Atmosphere
run 0–F and Law/Tech run 0–L, so every value the dice can produce is described.

---

## Starport (2D)

The referee selects the subsector's nature, then rolls on that column. This table
is identical to MegaTraveller's.

| Die | Backwater | Standard | Mature | Cluster |
| --- | --------- | -------- | ------ | ------- |
| 2   | A         | A        | A      | A       |
| 3   | A         | A        | A      | A       |
| 4   | B         | A        | A      | A       |
| 5   | B         | B        | B      | A       |
| 6   | C         | B        | B      | B       |
| 7   | C         | C        | C      | B       |
| 8   | C         | C        | C      | C       |
| 9   | D         | D        | D      | C       |
| 10  | E         | E        | E      | D       |
| 11  | E         | E        | E      | E       |
| 12  | X         | X        | E      | X       |

- **Backwater** — out of the mainstream of interstellar culture/communication.
- **Standard** — the expected norm.
- **Mature** — an older, more established system.
- **Cluster** — many worlds close together.

### Starport quality

| Type | Quality   | Shipyards  | Repair       | Fuel      |
| ---- | --------- | ---------- | ------------ | --------- |
| A    | Excellent | Starships  | Overhaul     | Refined   |
| B    | Good      | Spacecraft | Overhaul     | Refined   |
| C    | Routine   | —          | Major damage | Unrefined |
| D    | Poor      | —          | Minor damage | Unrefined |
| E    | Frontier  | —          | —            | —         |
| X    | None      | —          | —            | —         |

---

## Characteristic tables

Extended-hex throughout (0–9, A–Z omitting I and O).

### World Size (2D − 2)

The 2D−2 roll reaches 0–A. (The R / S / SGG / LGG codes in the source table are
for rings, small satellites, and gas giants — extended generation, out of scope.)

| Code | Description             | Diameter          |
| ---- | ----------------------- | ----------------- |
| 0    | Asteroid/Planetoid Belt | (multiple bodies) |
| 1    | Small                   | 1,600 km          |
| 2    | Small (Luna)            | 3,200 km          |
| 3    | Small (Mercury)         | 4,800 km          |
| 4    | Small (Mars)            | 6,400 km          |
| 5    | Medium                  | 8,000 km          |
| 6    | Medium                  | 9,600 km          |
| 7    | Medium                  | 11,200 km         |
| 8    | Large (Terra)           | 12,800 km         |
| 9    | Large                   | 14,400 km         |
| A    | Large                   | 16,000 km         |

### Atmosphere (2D − 7 + Size; Size 0 → 0)

| Code | Description                 | Pressure (atm) |
| ---- | --------------------------- | -------------- |
| 0    | Vacuum                      | 0.00           |
| 1    | Vacuum (trace)              | 0.001–0.09     |
| 2    | Vacuum (very thin, tainted) | 0.10–0.42      |
| 3    | Vacuum (very thin)          | 0.10–0.42      |
| 4    | Thin (tainted)              | 0.43–0.70      |
| 5    | Thin                        | 0.43–0.70      |
| 6    | Standard                    | 0.71–1.49      |
| 7    | Standard (tainted)          | 0.71–1.49      |
| 8    | Dense                       | 1.50–2.49      |
| 9    | Dense (tainted)             | 1.50–2.49      |
| A    | Exotic                      | Varies         |
| B    | Exotic (corrosive)          | Varies         |
| C    | Exotic (insidious)          | Varies         |
| D    | Exotic (dense, high)        | Varies         |
| E    | Exotic (thin, low)          | Varies         |
| F    | Exotic (ellipsoid)          | Varies         |

### Hydrographics (2D − 7 + Size; Size ≤ 1 → 0; atm ≤ 1 or ≥ A → DM −4)

| Code | Description  | % water |
| ---- | ------------ | ------- |
| 0    | Desert World | 0%      |
| 1    | Dry World    | 10%     |
| 2    | Dry World    | 20%     |
| 3    | Wet World    | 30%     |
| 4    | Wet World    | 40%     |
| 5    | Wet World    | 50%     |
| 6    | Wet World    | 60%     |
| 7    | Wet World    | 70%     |
| 8    | Wet World    | 80%     |
| 9    | Wet World    | 90%     |
| A    | Water World  | 100%    |

### Population (2D − 2)

Digit is the exponent of 10 (people ≈ 10^digit).

| Code | Description                 | Population      |
| ---- | --------------------------- | --------------- |
| 0    | Inconsequential (< ten)     | 0–9             |
| 1    | Inconsequential (tens)      | 10–99           |
| 2    | Inconsequential (hundreds)  | 100–999         |
| 3    | Low (thousands)             | 1,000–9,999     |
| 4    | Low (ten-thousands)         | 10,000–99,999   |
| 5    | Low (hundred-thousands)     | 100,000–999,999 |
| 6    | Moderate (millions)         | 1,000,000s      |
| 7    | Moderate (ten-millions)     | 10,000,000s     |
| 8    | Moderate (hundred-millions) | 100,000,000s    |
| 9    | High (billions)             | 1,000,000,000s  |
| A    | High (ten-billions)         | 10,000,000,000s |

### Government (2D − 7 + Population)

| Code | Description                 |
| ---- | --------------------------- |
| 0    | No Government Structure     |
| 1    | Company/Corporation         |
| 2    | Participating Democracy     |
| 3    | Self-Perpetuating Oligarchy |
| 4    | Representative Democracy    |
| 5    | Feudal Technocracy          |
| 6    | Captive Government / Colony |
| 7    | Balkanization               |
| 8    | Civil Service Bureaucracy   |
| 9    | Impersonal Bureaucracy      |
| A    | Charismatic Dictator        |
| B    | Non-Charismatic Leader      |
| C    | Charismatic Oligarchy       |
| D    | Religious Dictatorship      |
| E    | Religious Autocracy         |
| F    | Totalitarian Oligarchy      |

### Law Level (2D − 7 + Government)

Identical to MegaTraveller's table.

| Code | Description                                               |
| ---- | --------------------------------------------------------- |
| 0    | No law (no prohibitions)                                  |
| 1    | Low law (body pistols, explosives, poison gas prohibited) |
| 2    | Low law (portable energy weapons prohibited)              |
| 3    | Low law (machine guns, automatic rifles prohibited)       |
| 4    | Moderate law (light assault weapons prohibited)           |
| 5    | Moderate law (personal concealable weapons prohibited)    |
| 6    | Moderate law (all firearms except shotguns prohibited)    |
| 7    | Moderate law (shotguns prohibited)                        |
| 8    | High law (blade weapons controlled, no open display)      |
| 9    | High law (weapon possession outside the home prohibited)  |
| A    | Extreme law (weapon possession prohibited)                |
| B    | Extreme law (rigid control of civilian movement)          |
| C    | Extreme law (unrestricted invasion of privacy)            |
| D    | Extreme law (paramilitary law enforcement)                |
| E    | Extreme law (full-fledged police state)                   |
| F    | Extreme law (all facets of daily life rigidly controlled) |
| G    | Extreme law (severe punishment for petty infractions)     |
| H    | Extreme law (legalized oppressive practices)              |
| J    | Extreme law (routinely oppressive and restrictive)        |
| K    | Extreme law (excessively oppressive and restrictive)      |
| L    | Extreme law (totally oppressive and restrictive)          |

### Tech Level (1D + tech DMs)

| Code | Description                                | ~Historical  |
| ---- | ------------------------------------------ | ------------ |
| 0    | Pre-Industrial (primitive)                 | Stone Age    |
| 1    | Pre-Industrial (bronze, iron)              | Middle Ages  |
| 2    | Pre-Industrial (printing press)            | circa 1600   |
| 3    | Pre-Industrial (basic science)             | circa 1800   |
| 4    | Industrial (internal combustion)           | circa 1900   |
| 5    | Industrial (mass production)               | circa 1930   |
| 6    | Pre-Stellar (nuclear power)                | circa 1950   |
| 7    | Pre-Stellar (miniaturized electronics)     | circa 1970   |
| 8    | Pre-Stellar (superconductors)              | circa 1990   |
| 9    | Early Stellar (fusion power)               | circa 2010   |
| A    | Early Stellar (jump drive)                 | circa 2100   |
| B    | Average Stellar (large starships)          |              |
| C    | Average Stellar (sophisticated robots)     |              |
| D    | Average Stellar (holographic data storage) | the Imperium |
| E    | High Stellar (anti-grav cities)            |              |
| F    | High Stellar (anagathics)                  |              |
| G    | High Stellar (global terraforming)         |              |
| H    | Extreme Stellar                            |              |
| J    | Extreme Stellar                            |              |
| K    | Extreme Stellar                            |              |
| L    | Extreme Stellar                            | the Ancients |

---

## Technology DM matrix

Roll 1D, add every applicable DM, floor at 0. The maximum adjusted roll is **20**
(Starport A, World Size 1, Atmosphere 3, Population A, Government 5), which
reconciles at Starport A = **+6** — the value printed here.

| Value        | Starport | Size | Atmosphere | Hydrographics | Population | Government |
| ------------ | -------- | ---- | ---------- | ------------- | ---------- | ---------- |
| 0            | —        | +2   | +1         | —             | —          | +1         |
| 1            | —        | +2   | +1         | —             | +1         | —          |
| 2            | —        | +1   | +1         | —             | +1         | —          |
| 3            | —        | +1   | +1         | —             | +1         | —          |
| 4            | —        | +1   | —          | —             | +1         | —          |
| 5            | —        | —    | —          | —             | +1         | +1         |
| 6–8          | —        | —    | —          | —             | —          | —          |
| 9            | —        | —    | —          | +1            | +2         | —          |
| A            | +6       | —    | +1         | +2            | +4         | —          |
| B            | +4       | —    | +1         | —             | —          | —          |
| C            | +2       | —    | +1         | —             | —          | —          |
| D            | —        | —    | +1         | —             | —          | −2         |
| E            | —        | —    | +1         | —             | —          | −1         |
| F            | +1       | —    | +1         | —             | —          | −1         |
| X (starport) | −4       | —    | —          | —             | —          | —          |

---

## Bases

Rolled per starport (the "Base Presence" table); E and X starports never have
bases. Identical to MegaTraveller.

| Starport | Imperial Naval | Imperial Scout | Non-Imperial Military |
| -------- | -------------- | -------------- | --------------------- |
| A        | 8+             | 10+            | 10+                   |
| B        | 8+             | 9+             | 9+                    |
| C        | —              | 8+             | 8+                    |
| D        | —              | 7+             | —                     |
| E        | —              | —              | —                     |
| X        | —              | —              | —                     |

### Base codes

| Code | Meaning                                                                 |
| ---- | ----------------------------------------------------------------------- |
| N    | Imperial Naval base                                                     |
| S    | Imperial Scout base                                                     |
| A    | Imperial Naval **and** Scout base                                       |
| B    | Imperial Naval base **and** Scout Way Station (an extensive scout base) |
| M    | Non-Imperial Military base                                              |

> As in MegaTraveller, code `B` (Way Station) has no separate generation throw —
> a Way Station is simply an extensive scout base. The generator emits N/S/A/M
> and defers `B`.

---

## Out of scope

These are deferred physical/setting layers, not part of the mainworld UWP:

- **Planetary Density → surface gravity** (TNE core p. 190): step 4 also rolls
  1D on a Planetary Density table and derives surface gravity from size ×
  density. Not part of the seven UWP digits.
- **World Tamer's Handbook "Detailing Planets"** layer: temperature (from star +
  orbit + greenhouse), detailed atmospheric pressure/composition/taint, world
  mass, and mapping — all built on top of a base UWP.
- **Extended System Generation** (stars, orbits, satellites, gas giants,
  planetoid belts) and the Collapse Effects sequence.
