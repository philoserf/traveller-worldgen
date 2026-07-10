# Traveller5 (T5) — Mainworld Generation

Extracted from _Traveller5 Core Book 3: Worlds and Adventures_ (Far Future
Enterprises, print/PDF edition 5.10), the "World Generation" section: the
**WorldGen Checklist (A)** and **Basics (B)** on printed pp. 23–24, the
**StSAHPGL-T** UWP tables on p. 25, **Trade Classes (D)** on p. 26, the **Ix Ex
Cx** extensions on p. 27, and **WorldGen NABZ NIL (F)** on p. 28.

Scope: the **mainworld** plus the system-level counts and extensions that derive
from it — matching the `megatraveller` package and then some. Generated here:

- the UWP (**StSAHPGL-T**),
- Naval and Scout **bases**,
- **PBG** — Population digit, Planetoid **B**elts, **G**as giants,
- the generated **Trade Classifications**,
- the **Ix** (Importance), **Ex** (Economic), and **Cx** (Cultural) extensions.

> **Verified against the PDF (2026-07-10).** Cells transcribed from Core Book 3
> pp. 23–28. Where the source defers a step to the referee or to system-map
> context we don't model, it is noted below and omitted rather than invented.

## How T5 differs from the earlier editions

The other three editions in this repo (`classic`, `megatraveller`, `tne`) share
one core mechanic — `2D − 7 + characteristic`. **T5 does not.** Its core draw is
**Flux**.

- **Flux** `= 1D − 1D`, an integer in **−5…+5** (a signed die). Atmosphere,
  Government, Law, and the extensions are all Flux-based.
- **Hydrographics keys off Atmosphere**, not Size: `Hyd = Flux + Atm + mods`.
  (Every other edition uses `2D − 7 + Size`.)
- **Ranges are wider and are hard-capped**: Size 0–F, Atmosphere 0–F,
  Government 0–F, Law 0–J, Hydrographics 0–A. Because generation clamps to these
  caps, no characteristic is ever "beyond described range."
- **Starport is a single 2D table** (not the four-column subsector-nature table
  of Mega/TNE) and is **identical to Classic Book 3**.
- **Bases roll `2D ≤ target`** (T5's "N−" notation means "N or less"), the mirror
  of Classic's `2D ≥ target`.
- **Gas giants and planetoid belts** use flat formulas (`2D/2 − 2`, `1D − 3`),
  not Mega's presence-then-quantity rolls.

### Deferred (system-map / referee-assigned, out of scope)

Faithful to a bare mainworld generator, these documented T5 steps are **not**
generated:

- **MainWorld Type** (Planet vs. Satellite of a gas giant) — always treated as a
  planet; satellite orbits are a system-map concern.
- **Habitable Zone / Climate** — orbit-dependent. This is why the **Climate**
  trade codes (Fr, Ho, Co, Tr, Tu, Tz) are omitted (p. 26 itself notes "Lk, Tz,
  Ho, and Co refer to climate but are not properly TCs").
- **Secondary** trade codes (Fa, Mi, Mr, Pe, Re) — marked "not MW" / referee /
  HZ-dependent on p. 26.
- **Political** (Cp, Cs, Cx, Cy) and **Special** (Fo, Pz, Da, Ab, An) trade codes
  — "assigned by Referee (not generated)" per p. 26.
- **Nobility (N)**, **Allegiance (A)**, **Travel Zone (Z)**, **Native
  Intelligent Life (NIL)** — the rest of the NABZ-NIL page (p. 28); referee /
  setting data.
- **Depot** and **Way Station** bases — "1 per 1000 worlds" / "1 per 50 parsecs
  on a trade route" (p. 28); the generator emits only N/S. The Ix "+1 if Way
  Station" term is therefore always 0 here.
- The system map itself (stars, orbits, other worlds, satellites — checklist
  steps F–K and G).

### Tech Level descriptions

Core Book 3 p. 25 gives only the **TL DM formula**, not a prose table. Tech-level
_values_ are generated exactly per that formula; the era-band descriptions used
for display are the edition-invariant Traveller TL bands (identical to the `tne`
package's, since these labels do not change between editions). A generated TL
tops out at 20 — the +2 size DM (size 0–1) and +2 hydrographics DM (hydro A)
cannot both apply, since size < 2 forces hydrographics to 0 — so every reachable
value has a band and "(beyond described range)" never surfaces.

---

## Generation sequence (mainworld)

Adjusted rolls below their floor are treated as the floor; rolls above a cap are
set to the cap. "Flux" is one `1D − 1D` draw (two dice).

| Step | Characteristic    | Throw                                                      |
| ---- | ----------------- | ---------------------------------------------------------- |
| 1    | Starport          | 2D → [Starport table](#starport-s)                         |
| 2    | Naval base        | starport A/B only: 2D ≤ target                             |
| 3    | Scout base        | starport A/B/C/D only: 2D ≤ target                         |
| 4    | Size              | 2D − 2; **if 10, reroll 1D + 9** (→ A–F)                   |
| 5    | Atmosphere        | Flux + Size (Size 0 → 0); floor 0, cap F                   |
| 6    | Hydrographics     | Flux + Atm (−4 if Atm < 2 or Atm > 9); Size < 2 → 0; cap A |
| 7    | Population (exp.) | 2D − 2; **if 10, reroll 2D + 3** (→ A–F)                   |
| 8    | Government        | Flux + Pop; floor 0, cap F                                 |
| 9    | Law Level         | Flux + Gov; floor 0, cap J                                 |
| 10   | Tech Level        | 1D + [tech DMs](#tech-level-tl)                            |
| 11   | Population digit  | if Pop > 0: even 1–9 (else 0)                              |
| 12   | Planetoid belts   | 1D − 3, floor 0                                            |
| 13   | Gas giants        | 2D / 2 − 2, floor 0                                        |
| 14   | Trade classes     | derived from the UWP                                       |
| 15   | Ix / Ex / Cx      | derived from the UWP + trade + PBG (Ex/Cx draw more dice)  |
| 16   | Name              | invented (pronounceable syllables)                         |

**Draw order is fixed and load-bearing** (see the determinism contract in the
repo's CLAUDE.md). A conditional die is consumed **only** when its branch
applies: the Size/Pop rerolls, the Atmosphere/Hydrographics Flux (skipped for the
automatic Size 0 / Size < 2 results), the base rolls (skipped for starports that
can't hold that base), the Population-digit roll (skipped when Pop = 0), and the
Ex Infrastructure and Cx rolls (see those sections).

The Ex/Cx extra draws, in order: **Ex** — Resources (2D), Infrastructure (1D if
Pop 4–6, 2D if Pop ≥ 7, none if Pop ≤ 3), Efficiency (Flux); **Cx** (only if
Pop > 0) — Heterogeneity (Flux), Strangeness (Flux), Symbols (Flux).

---

## Starport (S)

Single 2D table (Core Book 3 p. 24, table 1A "Starports on the Mainworld").
Identical to Classic Book 3.

| 2D    | Type | Quality   |
| ----- | ---- | --------- |
| 2–4   | A    | Excellent |
| 5–6   | B    | Good      |
| 7–8   | C    | Routine   |
| 9     | D    | Poor      |
| 10–11 | E    | Frontier  |
| 12    | X    | None      |

### Bases

T5 "N−" means "present if 2D rolls N or less" (p. 24 table 1A, confirmed by p. 28
table B). Only starports A–D can hold a generated base; E and X never do.

| Base  | A    | B    | C    | D    |
| ----- | ---- | ---- | ---- | ---- |
| Naval | 2D≤6 | 2D≤5 | —    | —    |
| Scout | 2D≤4 | 2D≤5 | 2D≤6 | 2D≤7 |

Base codes: **N** (Naval), **S** (Scout), **A** (both). Depot and Way Station are
deferred (see scope).

---

## UWP — StSAHPGL-T (p. 25)

### Size (S)

`Size = 2D − 2`. If the result is 10, reroll `1D + 9` (giving A–F). Floor 0.

| Digit | Diameter             | Digit | Diameter              |
| ----- | -------------------- | ----- | --------------------- |
| 0     | Asteroid Belt        | 8     | 8,000 mi / 12,800 km  |
| 1     | 1,000 mi / 1,600 km  | 9     | 9,000 mi / 14,400 km  |
| 2     | 2,000 mi / 3,200 km  | A     | 10,000 mi / 16,000 km |
| 3     | 3,000 mi / 4,800 km  | B     | 11,000 mi / 17,600 km |
| 4     | 4,000 mi / 6,400 km  | C     | 12,000 mi / 19,200 km |
| 5     | 5,000 mi / 8,000 km  | D     | 13,000 mi / 20,800 km |
| 6     | 6,000 mi / 9,600 km  | E     | 14,000 mi / 22,400 km |
| 7     | 7,000 mi / 11,200 km | F     | 15,000 mi / 24,000 km |

### Atmosphere (A)

`Atm = Flux + Size`. If Atm < 0 or Size = 0, Atm = 0. If Atm > F, Atm = F.

| Digit | Description       | Digit | Description   |
| ----- | ----------------- | ----- | ------------- |
| 0     | Vacuum            | 8     | Dense         |
| 1     | Trace             | 9     | Dense Tainted |
| 2     | Very Thin Tainted | A     | Exotic        |
| 3     | Very Thin         | B     | Corrosive     |
| 4     | Thin Tainted      | C     | Insidious     |
| 5     | Thin              | D     | Dense, High   |
| 6     | Standard          | E     | Thin, Low     |
| 7     | Standard Tainted  | F     | Unusual       |

### Hydrographics (H)

`Hyd = Flux + Atm + mods`. If Size < 2, Hyd = 0. If Atm < 2 or Atm > 9, DM −4. If
Hyd < 0, Hyd = 0. If Hyd > A, Hyd = A.

| Digit | Description  | Digit | Description |
| ----- | ------------ | ----- | ----------- |
| 0     | Desert World | 6     | 60% Water   |
| 1     | 10% Water    | 7     | 70% Water   |
| 2     | 20% Water    | 8     | 80% Water   |
| 3     | 30% Water    | 9     | 90% Water   |
| 4     | 40% Water    | A     | Water World |
| 5     | 50% Water    |       |             |

### Population (P)

`Pop = 2D − 2`. If the result is 10, reroll `2D + 3` (giving A–F). The digit is
an exponent of 10.

| Digit | Description              | Digit | Description               |
| ----- | ------------------------ | ----- | ------------------------- |
| 0     | Unpopulated              | 8     | Hundred Millions (10^8)   |
| 1     | Tens (10^1)              | 9     | Billions (10^9)           |
| 2     | Hundreds (10^2)          | A     | Ten Billions (10^10)      |
| 3     | Thousands (10^3)         | B     | Hundred Billions (10^11)  |
| 4     | Ten Thousands (10^4)     | C     | Trillions (10^12)         |
| 5     | Hundred Thousands (10^5) | D     | Ten Trillions (10^13)     |
| 6     | Millions (10^6)          | E     | Hundred Trillions (10^14) |
| 7     | Ten Millions (10^7)      | F     | Quadrillions (10^15)      |

The separate **Population digit** (PBG's first digit) is an even 1–9 when Pop > 0,
else 0.

### Government (G)

`Gov = Flux + Pop`. Floor 0. If Gov > F, Gov = F.

| Digit | Description                 | Digit | Description                  |
| ----- | --------------------------- | ----- | ---------------------------- |
| 0     | No Government Structure     | 8     | Civil Service Bureaucracy    |
| 1     | Company / Corporation       | 9     | Impersonal Bureaucracy       |
| 2     | Participating Democracy     | A     | Charismatic Dictatorship     |
| 3     | Self-Perpetuating Oligarchy | B     | Non-Charismatic Dictatorship |
| 4     | Representative Democracy    | C     | Charismatic Oligarchy        |
| 5     | Feudal Technocracy          | D     | Religious Dictatorship       |
| 6     | Captive Government / Colony | E     | Religious Autocracy          |
| 7     | Balkanization               | F     | Totalitarian Oligarchy       |

### Law Level (L)

`Law = Flux + Gov`. Floor 0. If Law > J, Law = J.

| Digit | Description                                    |
| ----- | ---------------------------------------------- |
| 0     | No Law. No prohibitions.                       |
| 1     | Low Law. Prohibition of WMD, Psi weapons.      |
| 2     | Low Law. Prohibition of "Portable" Weapons.    |
| 3     | Low Law. Prohibition of Acid, Fire, Gas.       |
| 4     | Moderate Law. Prohibition of Laser, Beam.      |
| 5     | Moderate Law. No Shock, EMP, Rad, Mag, Grav.   |
| 6     | Moderate Law. Prohibition of MachineGuns.      |
| 7     | Moderate Law. Prohibition of Pistols.          |
| 8     | High Law. Open display of weapons prohibited.  |
| 9     | High Law. No weapons outside the home.         |
| A     | Extreme Law. All weapons prohibited.           |
| B     | Extreme Law. Continental passports required.   |
| C     | Extreme Law. Unrestricted invasion of privacy. |
| D     | Extreme Law. Paramilitary law enforcement.     |
| E     | Extreme Law. Full-fledged police state.        |
| F     | Extreme Law. Daily life rigidly controlled.    |
| G     | Extreme Law. Disproportionate punishment.      |
| H     | Extreme Law. Legalized oppressive practices.   |
| J     | Extreme Law. Routine oppression.               |

### Tech Level (TL)

`TL = 1D + DMs`, floor 0. DMs (p. 25):

| Source        | DM  | Source       | DM  |
| ------------- | --- | ------------ | --- |
| Starport A    | +6  | Hyd 9        | +1  |
| Starport B    | +4  | Hyd A        | +2  |
| Starport C    | +2  | Pop 1–5      | +1  |
| Starport X    | −4  | Pop 9        | +2  |
| Size 0–1      | +2  | Pop A+ (≥10) | +4  |
| Size 2–4      | +1  | Gov 0        | +1  |
| Atm 0–3       | +1  | Gov 5        | +1  |
| Atm A–F (≥10) | +1  | Gov D (13)   | −2  |

(The source also lists "Spaceport F +1", which applies to non-mainworld
spaceports F/G/H and never to a mainworld starport A–X.)

---

## Trade Classifications (p. 26)

A world carries every code whose characteristics all fall in the listed digit
sets ("—" = unconstrained). Only the **Planetary**, **Population**, and
**Economic** groups are generated here; Climate, Secondary, Political, and
Special are deferred (see scope). Order follows the source table.

| Code | Size | Atm             | Hyd | Pop | Gov | Law | Meaning           |
| ---- | ---- | --------------- | --- | --- | --- | --- | ----------------- |
| As   | 0    | 0               | 0   | —   | —   | —   | Asteroid Belt     |
| De   | —    | 2–9             | 0   | —   | —   | —   | Desert            |
| Fl   | —    | A–C             | 1–A | —   | —   | —   | Fluid             |
| Ga   | 6–8  | 5,6,8           | 5–7 | —   | —   | —   | Garden World      |
| He   | 3–C  | 2,4,7,9,A–C     | 0–2 | —   | —   | —   | Hellworld         |
| Ic   | —    | 0,1             | 1–A | —   | —   | —   | Ice-Capped        |
| Oc   | A–F  | 3–9,D–F         | A   | —   | —   | —   | Ocean World       |
| Va   | —    | 0               | —   | —   | —   | —   | Vacuum            |
| Wa   | 3–9  | 3–9,D–F         | A   | —   | —   | —   | Water World       |
| Di   | —    | —               | —   | 0   | 0   | 0   | Dieback (TL ≥ 1)  |
| Ba   | —    | —               | —   | 0   | 0   | 0   | Barren (TL 0)     |
| Lo   | —    | —               | —   | 1–3 | —   | —   | Low Population    |
| Ni   | —    | —               | —   | 4–6 | —   | —   | Non-Industrial    |
| Ph   | —    | —               | —   | 8   | —   | —   | Pre-High Pop.     |
| Hi   | —    | —               | —   | 9–F | —   | —   | High Population   |
| Pa   | —    | 4–9             | 4–8 | 4,8 | —   | —   | Pre-Agricultural  |
| Ag   | —    | 4–9             | 4–8 | 5–7 | —   | —   | Agricultural      |
| Na   | —    | 0–3             | 0–3 | 6–F | —   | —   | Non-Agricultural  |
| Px   | —    | 2,3,A,B         | 1–5 | 3–6 | —   | 6–9 | Prison/Exile Camp |
| Pi   | —    | 0,1,2,4,7,9     | —   | 7,8 | —   | —   | Pre-Industrial    |
| In   | —    | 0,1,2,4,7,9,A–C | —   | 9–F | —   | —   | Industrial        |
| Po   | —    | 2–5             | 0–3 | —   | —   | —   | Poor              |
| Pr   | —    | 6,8             | —   | 5,9 | —   | —   | Pre-Rich          |
| Ri   | —    | 6,8             | —   | 6–8 | —   | —   | Rich              |

Di and Ba share the same S/A/H/P/G/L constraints (Pop 0, Gov 0, Law 0); they are
distinguished by Tech Level — **Di** (Dieback) has TL ≥ 1 (technology remains,
population gone), **Ba** (Barren) has TL 0 — per the "Dieback (000-T)" comment.

Unlike MegaTraveller, T5's table states no As→Va suppression rule, so an asteroid
belt legitimately carries **both** As and Va; both codes are emitted.

---

## Extensions — Ix Ex Cx (p. 27)

### Importance Extension (Ix), rendered `{±n}`

Sum of:

| Condition                      | Value                    |
| ------------------------------ | ------------------------ |
| Starport A or B                | +1                       |
| Starport D or worse (D/E/X)    | −1                       |
| Tech Level ≥ G (16)            | +1                       |
| Tech Level ≥ A (10)            | +1                       |
| Tech Level ≤ 8                 | −1                       |
| Each of Ag, Hi, In, Ri present | +1 each                  |
| Population ≤ 6                 | −1                       |
| Naval **and** Scout base       | +1                       |
| Way Station                    | +1 (always 0 — deferred) |

(The three TL terms stack: TL ≥ 16 scores both the G and A rows.)

### Economic Extension (Ex), rendered `(RLI±E)`

- **Resources** `= 2D`; if TL ≥ 8, add Gas Giants + Planetoid Belts. Floor 0.
- **Labor** `= Pop − 1`. Floor 0.
- **Infrastructure** `=` 0 if Pop = 0; Ix if Pop 1–3; `1D + Ix` if Pop 4–6;
  `2D + Ix` if Pop ≥ 7. Floor 0.
- **Efficiency** `= Flux`. May be negative (rendered with sign).

**RU** (Resource Units) `= Resources × Labor × Infrastructure × Efficiency`;
any factor of 0 is treated as 1 (to avoid zeroing the product), so RU can be
negative when Efficiency is.

### Cultural Extension (Cx), rendered `[HASS]`

If Pop = 0, all four are 0. Otherwise each is floored at 1 ("less than 1 = 1"):

| Component     | Value      |
| ------------- | ---------- |
| Heterogeneity | Pop + Flux |
| Acceptance    | Pop + Ix   |
| Strangeness   | Flux + 5   |
| Symbols       | Flux + TL  |

R, L, I and the Cx components are single characters rendered in extended hex.
