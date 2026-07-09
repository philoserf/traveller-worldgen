# MegaTraveller — Basic Mainworld Generation

Extracted from _MegaTraveller Referee's Manual_ (Digest Group Publications / GDW,
1987), "Basic Mainworld Generation" flowchart (printed pp. 24–25), with UWP code
definitions from the "Universal World Profile Tables" (printed pp. 22–23) and the
_Imperial Encyclopedia_.

Scope: the **mainworld** only (single world), matching the `classic` package.
Extended System Generation (orbits, satellites, additional planets) is out of
scope.

> **Verified against the PDF (2026-07-09).** The two previously flagged cells,
> plus the Trade Classification ranges, were read directly from the scanned
> _Referee's Manual_ (pp. 24–25) and _Imperial Encyclopedia_ (p. 48):
>
> 1. **Technology DM matrix** — matches the reconstruction except that the
>    printed table cell for **Starport A reads `+8`**. This contradicts the same
>    page's worked "maximum adjusted roll is 20" example, which reconciles only
>    at `+6` (and matches Classic). We treat the `+8` cell as a print error and
>    use **`+6`**. See § Technology DM matrix.
> 2. **Way Station (`B`) base code** — there is _no_ generation throw for it.
>    The _Referee's Manual_ Basic Mainworld Generation has no Bases step at all,
>    and the _Imperial Encyclopedia_ describes a Way Station only narratively (a
>    large scout base servicing ships ≤10,000 tons, "the equivalent of a naval
>    base"). Code `B` stays deferred; the generator emits N/S/A/M. See § Bases.

---

## How MegaTraveller differs from Classic Book 3

- **Same core formulas** — Size `2D−2`, Atmosphere `2D−7+Size`, Hydrographics
  `2D−7+Size` (−4 DM if atm ≤1 or ≥A), Population `2D−2`, Government
  `2D−7+Pop`, Law `2D−7+Gov`, Tech `1D+DMs`.
- **Starport** is a **4-column** table; the referee picks the subsector's nature
  (Backwater / Standard / Mature / Cluster).
- **Extended, complete ranges** — Government 0–F, Law 0–L, Atmosphere 0–F, Tech
  0–L are all fully described, so no value the dice can roll is "out of range."
- **Bases** add a Non-Imperial Military base and combined codes (N/S/A/B/M).
- **Trade classifications** — a derived layer (Ag, As, Ba, …) computed from the
  UWP. _New._
- **Gas giants** — rolled presence + quantity. _New._

---

## Generation sequence (mainworld)

| Step | Characteristic        | Throw                                                               |
| ---- | --------------------- | ------------------------------------------------------------------- |
| 1    | Starport              | 2D → [Starport table](#starport-2d) for the chosen subsector nature |
| 2    | World Size            | 2D − 2                                                              |
| 3    | Atmosphere            | 2D − 7 + Size (Size 0 → Atmosphere 0)                               |
| 4    | Hydrographics         | 2D − 7 + Size (Size ≤ 1 → 0; Atmosphere ≤ 1 or ≥ A → DM −4)         |
| 5    | Population            | 2D − 2                                                              |
| 6    | Government            | 2D − 7 + Population                                                 |
| 7    | Law Level             | 2D − 7 + Government                                                 |
| 8    | Tech Level            | 1D + [tech DMs](#technology-dm-matrix)                              |
| 9    | Bases                 | per-starport throws → N/S/A/B/M                                     |
| 10   | Trade classifications | derived from the UWP                                                |
| 11   | Gas giants            | 5+ on 2D for presence, then a quantity roll                         |
| 12   | Planetoid belts       | 8+ on 2D for presence, then a quantity roll                         |

Adjusted rolls below 0 are treated as 0. Because Gov/Law/Atmo/Tech maxima
(F/L/F/L) equal the highest the dice yield, every result has a defined
description.

The table lists Bases at step 9 because a base depends only on the starport, but
the **implementation draws the base rolls immediately after the starport**
(before World Size), matching the `classic` package. The base rolls need only the
starport, so the result is identical either way; rolling them early just keeps a
single, stable dice-draw order across editions for seed reproducibility. The
authoritative draw order is documented on `megatraveller.Generate`.

---

## Starport (2D)

The referee selects the subsector's nature, then rolls on that column.

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

| Code | Description                 |
| ---- | --------------------------- |
| 0    | Vacuum                      |
| 1    | Vacuum (trace)              |
| 2    | Vacuum (very thin, tainted) |
| 3    | Vacuum (very thin)          |
| 4    | Thin (tainted)              |
| 5    | Thin                        |
| 6    | Standard                    |
| 7    | Standard (tainted)          |
| 8    | Dense                       |
| 9    | Dense (tainted)             |
| A    | Exotic                      |
| B    | Corrosive                   |
| C    | Insidious                   |
| D    | Dense, high                 |
| E    | Ellipsoid                   |
| F    | Thin, low                   |

### Hydrographics (2D − 7 + Size; Size ≤ 1 → 0; atm ≤ 1 or ≥ A → DM −4)

| Code | Description  | % water |
| ---- | ------------ | ------- |
| 0    | Desert World | 0–4%    |
| 1    | Dry World    | 5–14%   |
| 2    | Dry World    | 15–24%  |
| 3    | Wet World    | 25–34%  |
| 4    | Wet World    | 35–44%  |
| 5    | Wet World    | 45–54%  |
| 6    | Wet World    | 55–64%  |
| 7    | Wet World    | 65–74%  |
| 8    | Wet World    | 75–84%  |
| 9    | Wet World    | 85–94%  |
| A    | Water World  | 95–100% |

### Population (2D − 2)

Digit is the exponent of 10 (people ≈ 10^digit).

| Code | Description                  | Population      |
| ---- | ---------------------------- | --------------- |
| 0    | Low (< ten)                  | 0–9             |
| 1    | Low (tens)                   | 10s             |
| 2    | Low (hundreds)               | 100s            |
| 3    | Low (thousands)              | 1,000s          |
| 4    | Moderate (ten-thousands)     | 10,000s         |
| 5    | Moderate (hundred-thousands) | 100,000s        |
| 6    | Moderate (millions)          | 1,000,000s      |
| 7    | Moderate (ten-millions)      | 10,000,000s     |
| 8    | Moderate (hundred-millions)  | 100,000,000s    |
| 9    | High (billions)              | 1,000,000,000s  |
| A    | High (ten-billions)          | 10,000,000,000s |

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
| 0    | Pre-Industrial (primitive)                 | stone age    |
| 1    | Pre-Industrial (bronze, iron)              | middle ages  |
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
| D    | Average Stellar (holographic data storage) |              |
| E    | High Stellar (anti-grav cities)            |              |
| F    | High Stellar (anagathics)                  |              |
| G    | High Stellar (global terraforming)         |              |
| H    | Extreme Stellar                            | the Imperium |
| J    | Extreme Stellar                            |              |
| K    | Extreme Stellar                            |              |
| L    | Extreme Stellar                            | the Ancients |

---

## Technology DM matrix

> **Verified against the PDF.** Read cell-by-cell from the printed grid; every
> value below matches the scan **except the Starport A cell, which prints `+8`**.
> That `+8` contradicts the same page's stated maximum adjusted roll of **20**:
> in the max example (Starport A, Size 1, Atmo 3, Pop A, Gov 5) the non-starport
> DMs sum to +8, so `1D(6) + starport + 8 = 20` requires starport `+6`. We use
> **`+6`** (the Classic value) and treat the printed `+8` as a print error. Roll
> 1D, add every applicable DM, floor at 0.

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
| E            | —        | —    | +1         | —             | —          | —          |
| F            | —        | —    | +1         | —             | —          | —          |
| X (starport) | −4       | —    | —          | —             | —          | —          |

---

## Bases

Rolled per starport; E and X starports never have bases.

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

> **Verified against the PDF.** The Basic Mainworld Generation flowchart has **no
> Bases step** (its steps are System Presence, System Details, Starport, Size,
> Atmosphere, Hydrographics, Population, Government, Law, Tech, Trade,
> Supplemental Remarks, Population Multiplier, Gas Giants, Planetoid Belts, Travel
> Zones, Allegiance — none assigns bases). The throw table above therefore
> follows standard/Classic Traveller convention rather than a printed MT
> generation table. The _Imperial Encyclopedia_ (p. 48) describes a **Way
> Station** only narratively — a large scout base for express-boat overhaul,
> "the equivalent of a naval base," servicing ships ≤10,000 tons — with no
> promotion throw. Code `B` therefore stays deferred; the generator emits
> N/S/A/M.

---

## Trade classifications (derived from the UWP)

A world can carry several. Ranges are inclusive; `—` means the characteristic is
unconstrained for that code.

| Code | Size | Atmo    | Hydro | Pop | Gov | Law | Meaning                  |
| ---- | ---- | ------- | ----- | --- | --- | --- | ------------------------ |
| Ag   | —    | 4–9     | 4–8   | 5–7 | —   | —   | Agricultural             |
| As   | 0    | 0       | 0     | —   | —   | —   | Asteroid                 |
| Ba   | —    | —       | —     | 0   | 0   | 0   | Barren                   |
| De   | —    | 2+      | 0     | —   | —   | —   | Desert                   |
| Fl   | —    | A+      | 1+    | —   | —   | —   | Fluid (non-water oceans) |
| Hi   | —    | —       | —     | 9+  | —   | —   | High Population          |
| Ic   | —    | 0–1     | 1+    | —   | —   | —   | Ice-Capped               |
| In   | —    | 2–4,7,9 | —     | 9+  | —   | —   | Industrial               |
| Lo   | —    | —       | —     | 0–3 | —   | —   | Low Population           |
| Na   | —    | 0–3     | 0–3   | 6+  | —   | —   | Nonagricultural          |
| Ni   | —    | —       | —     | 0–6 | —   | —   | Nonindustrial            |
| Po   | —    | 2–5     | 0–3   | —   | —   | —   | Poor                     |
| Ri   | —    | 6,8     | —     | 6–8 | 4–9 | —   | Rich                     |
| Va   | —    | 0       | —     | —   | —   | —   | Vacuum                   |
| Wa   | —    | —       | A     | —   | —   | —   | Water World              |

- An **Asteroid (As)** is automatically a **Vacuum (Va)** world (the Va code is
  then redundant).
- Race-specific notes (Aslan always Rich; Vargr disqualified from Rich if
  Government 7) are setting-dependent and out of scope for a bare generator.
- **Read from the PDF (corrections to an earlier reconstruction):** Industrial
  (In) atmosphere is `2–4,7,9` (= 2,3,4,7,9), and Rich (Ri) atmosphere is `6,8`
  (excluding tainted 7). The printed Fluid (Fl) row places `A+`/`1+` one column
  left (under Size/Atmo); this is a print misalignment — the intended and
  universally-used definition is Atmosphere `A+`, Hydrographics `1+`, kept above.

---

## Gas giants

Roll **5+ on 2D**; if present, roll 2D for the quantity.

| 2D  | Gas giants |
| --- | ---------- |
| 2   | 1          |
| 3   | 1          |
| 4   | 2          |
| 5   | 2          |
| 6   | 3          |
| 7   | 3          |
| 8   | 4          |
| 9   | 4          |
| 10  | 4          |
| 11  | 5          |
| 12  | 5          |

---

## Planetoid belts

Roll **8+ on 2D**; if present, roll 2D for the quantity (Referee's Manual step
16, p. 25).

| 2D  | Planetoid belts |
| --- | --------------- |
| 2   | 1               |
| 3   | 1               |
| 4   | 1               |
| 5   | 1               |
| 6   | 1               |
| 7   | 1               |
| 8   | 2               |
| 9   | 2               |
| 10  | 2               |
| 11  | 2               |
| 12  | 2               |
| 13  | 3               |

The printed table includes a **13** row (→ 3), but the step gives no DM, so a
plain 2D quantity roll reaches only 2–12 (at most 2 belts); the 13 row is
transcribed for fidelity but is unreachable here.

---

## Out of scope

Extended System Generation (orbits, satellites, additional planets), Population
Multiplier, Travel Zones, Supplemental Remarks, and Allegiance — referee/setting
layers or the deeper system model, deferred per project scope. (Gas giants and
planetoid belts, the two system-level counts on the printed flowchart, are
generated; the fuller orbit/star/body model is not.)
