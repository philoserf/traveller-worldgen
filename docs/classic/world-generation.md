# Classic Traveller — World Generation

Extracted from _Book 3: Worlds and Adventures_ (Game Designers' Workshop, 1977).

A world's Universal World Profile (UWP) is a string of seven characteristics plus a
technological index, expressed as single digits/letters (`A`–`Z`, omitting `I` and `O`
to avoid confusion with `1` and `0`). Illustrative profile (invented values —
the source gives no canonical example):

```
Starport Size Atmosphere Hydrographics Population Government LawLevel  TechLevel
   A       6       5           4            8          4         5    -   9
```

---

## Generation sequence (checklist)

1. **Map the subsector.** Throw 1D per hex; **4, 5, or 6** → a world is present.
2. **Starport type.** 2D, consult the [Starport table](#starport-type-2d).
3. **Space lanes.** Check all possible jump routes.
4. Generate each world's characteristics (all 2D unless noted):

| Step | Characteristic      | Throw                                                               |
| ---- | ------------------- | ------------------------------------------------------------------- |
| A    | Starport type       | (already rolled in step 2)                                          |
| B    | Planetary Size      | 2D − 2                                                              |
| C    | Atmosphere          | 2D − 7 + Size (Size 0 → Atmosphere 0)                               |
| D    | Hydrographics       | 2D − 7 + Size (Size 0 or 1 → Hyd 0; Atmosphere 0, 1, or A+ → DM −4) |
| E    | Population          | 2D − 2                                                              |
| F    | Government          | 2D − 7 + Population                                                 |
| G    | Law Level           | 2D − 7 + Government                                                 |
| H    | Technological Index | 1D + DMs from the [Tech Index Matrix](#technological-index)         |

5. Name each world and note its hex location.

> Base 2D rolls fall in 0–10, but DMs — including the size-based DMs on Atmosphere,
> Hydrographics, Government, and Law Level — can push results higher. Values above 9 are
> written as letters (A = 10, B = 11, C = 12, D = 13, …). Negative results are treated as 0.

---

## Star Mapping

### World occurrence

Throw 1D per hex; **4, 5, 6** = world present. Referee may apply a subsector-wide DM
of +1 or −1 to make worlds more/less frequent.

### Starport type (2D)

| Die (2D) | Type |
| -------- | ---- |
| 2        | A    |
| 3        | A    |
| 4        | A    |
| 5        | B    |
| 6        | B    |
| 7        | C    |
| 8        | C    |
| 9        | D    |
| 10       | E    |
| 11       | E    |
| 12       | X    |

### Jump routes

For each world, throw 1D and cross-reference the world pair's starport relationship
against the jump distance. If the die ≥ the tabled number, a space lane exists (draw a
line connecting the worlds). Examine each pair only once.

| World Pair | Jump-1 | Jump-2 | Jump-3 | Jump-4 |
| ---------- | ------ | ------ | ------ | ------ |
| A-A        | 1      | 2      | 4      | 5      |
| A-B        | 1      | 3      | 4      | 5      |
| A-C        | 1      | 4      | 6      | —      |
| A-D        | 1      | 5      | —      | —      |
| A-E        | 2      | —      | —      | —      |
| B-B        | 1      | 3      | 4      | 6      |
| B-C        | 2      | 4      | 6      | —      |
| B-D        | 3      | 6      | —      | —      |
| B-E        | 4      | —      | —      | —      |
| C-C        | 3      | 6      | —      | —      |
| C-D        | 4      | —      | —      | —      |
| C-E        | 4      | —      | —      | —      |
| D-D        | 4      | —      | —      | —      |
| D-E        | 5      | —      | —      | —      |
| E-E        | 6      | —      | —      | —      |

(1 hex = 1 parsec = 3.26 light years.)

---

## Characteristic tables

### Starport type

| Type | Description                                                                                                                                |
| ---- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| A    | Excellent quality. Refined fuel. Annual maintenance overhaul. Shipyard (starships and non-starships). Naval base on 8+; scout base on 10+. |
| B    | Good quality. Refined fuel. Annual maintenance overhaul. Shipyard (non-starships). Naval base on 8+; scout base on 9+.                     |
| C    | Routine quality. Unrefined fuel. Reasonable repair facilities. Scout base on 8+.                                                           |
| D    | Poor quality. Unrefined fuel. No repair or shipyard. Scout base on 7+.                                                                     |
| E    | Frontier installation. Bare bedrock — no fuel, facilities, or bases.                                                                       |
| X    | No starport. No provision for starship landings.                                                                                           |

### Planetary Size (2D − 2)

| Digit | Description                |
| ----- | -------------------------- |
| 0     | Asteroid/Planetoid Complex |
| 1     | 1,000 miles diameter       |
| 2     | 2,000 miles diameter       |
| 3     | 3,000 miles diameter       |
| 4     | 4,000 miles diameter       |
| 5     | 5,000 miles diameter       |
| 6     | 6,000 miles diameter       |
| 7     | 7,000 miles diameter       |
| 8     | 8,000 miles diameter       |
| 9     | 9,000 miles diameter       |
| A     | 10,000 miles diameter      |
| B     | 11,000 miles diameter      |
| C     | 12,000 miles diameter      |

### Planetary Atmosphere (2D − 7 + Size)

| Digit | Description        |
| ----- | ------------------ |
| 0     | No atmosphere      |
| 1     | Trace              |
| 2     | Very thin, tainted |
| 3     | Very thin          |
| 4     | Thin, tainted      |
| 5     | Thin               |
| 6     | Standard           |
| 7     | Standard, tainted  |
| 8     | Dense              |
| 9     | Dense, tainted     |
| A     | Exotic             |
| B     | Corrosive          |
| C     | Insidious          |

**Atmosphere notes:**

- **No atmosphere / Trace** — approximate vacuum; vacc suits required at all times.
- **Very thin** — respirator/compressor required for sufficient oxygen.
- **Thin / Standard / Dense** — breathable without assistance.
- **Tainted** — polluted; filter masks required in addition to any other apparatus.
- **Exotic** — oxygen tanks required, but protective suits are not.
- **Corrosive** — protective suits/measures similar to vacuum required.
- **Insidious** — like corrosive, but defeats personal protective measures in 2–12 hours.

### Hydrographic Percentage (2D − 7 + Size)

DM +Size applies; Size 0 or 1 → automatic 0; Atmosphere 0, 1, or A+ → DM −4.

| Digit | Description                |
| ----- | -------------------------- |
| 0     | No free standing water     |
| 1     | 10%                        |
| 2     | 20%                        |
| 3     | 30%                        |
| 4     | 40%                        |
| 5     | 50%                        |
| 6     | 60%                        |
| 7     | 70%                        |
| 8     | 80%                        |
| 9     | 90%                        |
| A     | All water. No land masses. |

### Population (2D − 2)

Population digit is an exponent of 10 (number of zeros after a 1).

| Digit | Description       |
| ----- | ----------------- |
| 0     | 0. No inhabitants |
| 1     | 10                |
| 2     | 100               |
| 3     | 1,000             |
| 4     | 10,000            |
| 5     | 100,000           |
| 6     | 1,000,000         |
| 7     | 10,000,000        |
| 8     | 100,000,000       |
| 9     | 1,000,000,000     |
| A     | 10,000,000,000    |

### Government (2D − 7 + Population)

| Type | Description                                                                                                               |
| ---- | ------------------------------------------------------------------------------------------------------------------------- |
| 0    | No government structure. Family bonds predominate.                                                                        |
| 1    | Company/Corporation. Managerial elite; citizens are employees/dependents.                                                 |
| 2    | Participating Democracy. Decisions by advice and consent of citizenry directly.                                           |
| 3    | Self-Perpetuating Oligarchy. Ruled by a restricted minority.                                                              |
| 4    | Representative Democracy. Elected representatives.                                                                        |
| 5    | Feudal Technocracy. Rule by those performing mutually beneficial technical activities.                                    |
| 6    | Captive Government. Imposed leadership answerable to an outside group; colony/conquered area.                             |
| 7    | Balkanization. No central authority; rival governments compete. Use the law level of the government nearest the starport. |
| 8    | Civil Service Bureaucracy. Agencies staffed by expertise-selected individuals.                                            |
| 9    | Impersonal Bureaucracy. Agencies insulated from the governed.                                                             |
| A    | Charismatic Dictator. Single leader with overwhelming confidence of citizens.                                             |
| B    | Non-Charismatic Leader. Successor to a charismatic dictator, through normal channels.                                     |
| C    | Charismatic Oligarchy. Select group enjoying overwhelming confidence of citizenry.                                        |
| D    | Religious Dictatorship. Religious organization ruling to its own needs.                                                   |

### Law Level (2D − 7 + Government)

Each law level includes all prohibitions of lower levels (thus shotguns,
prohibited at level 7, remain prohibited at every higher level). The described
scale ends at 9; a rolled law level above 9 keeps level 9's prohibitions while
the enforcement throw (below) continues to rise.

| Level | Prohibitions                                                                                               |
| ----- | ---------------------------------------------------------------------------------------------------------- |
| 0     | No weapons laws.                                                                                           |
| 1     | Body pistols (undetectable), explosives (bombs/grenades), poison gas prohibited.                           |
| 2     | Portable energy weapons (laser rifles/carbines) prohibited. Ship's gunnery unaffected.                     |
| 3     | Strict military weapons (machine guns, automatic rifles; not SMGs) prohibited.                             |
| 4     | Light assault weapons (submachine guns) prohibited.                                                        |
| 5     | Personal concealable firearms (pistols, revolvers) prohibited.                                             |
| 6     | Most firearms (all except shotguns) prohibited; open carry discouraged.                                    |
| 7     | Shotguns prohibited.                                                                                       |
| 8     | Long blades (all except daggers) controlled; open public possession prohibited (ownership not restricted). |
| 9     | Possession of any weapon outside one's home prohibited.                                                    |

> Law level also gives the throw for enforcement harassment: a person on a law-level-4
> world needs a saving throw of 4+ to avoid arrest when encountering an enforcement
> agent. Law does not apply to persons/ships at a starport (extraterritorial).

---

## Technological Index

Throw 1D and add DMs read from the matrix for each characteristic; sum for one total DM.

### Tech Index Matrix (DMs)

| Value | Starport | Size | Atm | Hyd | Pop | Govt |
| ----- | -------- | ---- | --- | --- | --- | ---- |
| 0     | —        | +2   | +1  | —   | —   | +1   |
| 1     | —        | +2   | +1  | —   | +1  | —    |
| 2     | —        | +1   | +1  | —   | +1  | —    |
| 3     | —        | +1   | +1  | —   | +1  | —    |
| 4     | —        | +1   | —   | —   | +1  | —    |
| 5     | —        | —    | —   | —   | +1  | +1   |
| 6     | —        | —    | —   | —   | —   | —    |
| 7     | —        | —    | —   | —   | —   | —    |
| 8     | —        | —    | —   | —   | —   | —    |
| 9     | —        | —    | —   | +1  | +2  | —    |
| A     | +6       | —    | +1  | +2  | +4  | —    |
| B     | +4       | —    | +1  | —   | —   | —    |
| C     | +2       | —    | +1  | —   | —   | —    |
| D     | —        | —    | +1  | —   | —   | −2   |
| E     | —        | —    | +1  | —   | —   | —    |
| X     | −4       | —    | —   | —   | —   | —    |

> Tech index ranges 0–18, commonly 4–10. Higher = greater capability.

### Technological Levels — Weapons, Computers, Communication

| TL  | Personal                      | Armor        | Special                    | Computers               | Communication |
| --- | ----------------------------- | ------------ | -------------------------- | ----------------------- | ------------- |
| 0   | Club, cudgel, Spear           |              |                            |                         | Runners       |
| 1   | Dagger, pike, Sword           | Jack         | Catapult                   | Abacus                  | Heliograph    |
| 2   | Halberd, Broadsword           |              | Cannon                     |                         |               |
| 3   | Foil, cutlass, Blade, bayonet |              |                            |                         |               |
| 4   | Revolver, Shotgun             | Cloth        | Artillery                  | Adding Machine          | Telephones    |
| 5   | Carbine, Rifle, Pistol, SMG   |              | Sandcasters, Mortars       | Model/1                 | Radio         |
| 6   | Auto Rifle                    |              | Missiles, Rocket Launchers | Model/1 bis             | Television    |
| 7   | Body Pistol                   | Mesh         | Pulse Laser                | Model/2                 |               |
| 8   | Laser Carbine                 |              | Auto-Cannon                | Model/2 bis             |               |
| 9   | Laser Rifle                   | Ablat        | Beam Laser                 | Model/3                 |               |
| 10  |                               | Reflec       |                            | Model/4                 |               |
| 11  |                               |              |                            | Model/5                 |               |
| 12  |                               |              |                            | Model/6                 |               |
| 13  |                               | Battle Dress |                            | Model/7                 |               |
| 17  |                               |              |                            | Artificial Intelligence |               |

### Technological Levels — Transportation & Fuels

| TL  | Water            | Land        | Air                  | Space            | Fuels       |
| --- | ---------------- | ----------- | -------------------- | ---------------- | ----------- |
| 0   | Canoes           | Carts       |                      |                  | Muscle      |
| 1   | Galley           | Wagons      |                      |                  |             |
| 2   |                  |             |                      |                  | Wind        |
| 3   | Sailing Ships    |             | Hot air balloon      |                  | Water wheel |
| 4   | Steamships       | Trains      | Dirigibles           |                  | Coal        |
| 5   |                  | Ground cars | Fixed wing aircraft  |                  | Oil         |
| 6   | Submersibles     | ATV, AFV    | Rotary wing aircraft |                  | Fission     |
| 7   | Hovercraft       | Hovercraft  |                      | Non-starships    | Solar       |
| 8   |                  |             | Air/Raft             |                  | Fusion      |
| 9   |                  |             |                      | Starships        |             |
| 10  |                  |             |                      | Drives H or less |             |
| 11  |                  |             |                      | Drives K or less |             |
| 12  |                  |             | Grav belts           | Drives N or less |             |
| 13  |                  |             |                      | Drives O or less |             |
| 14  |                  |             |                      | Drives U or less |             |
| 15  |                  |             |                      | All drives       |             |
| 16  | Matter Transport |             |                      |                  |             |
| 17  |                  |             |                      |                  | Anti-Matter |
| 18  |                  |             |                      |                  |             |
