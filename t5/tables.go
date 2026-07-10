package t5

// All tables in this file are transcribed from docs/t5/world-generation.md,
// which was read cell-by-cell against the scanned Traveller5 Core Book 3: Worlds
// and Adventures (Far Future Enterprises, edition 5.10, pp. 23–28). Descriptions
// are condensed from the source prose.

// beyondRange is the description returned for a characteristic value with no
// table entry. T5's UWP ranges (Size/Atm/Gov 0–F, Law 0–J, Hyd 0–A) and the
// reachable Tech Levels (0–20) are all complete, so this never surfaces in
// practice; it is a safety net.
const beyondRange = "(beyond described range)"

// starportDesc maps a starport code to its description (Book 3 p. 24, table 1A).
var starportDesc = map[byte]string{
	'A': "Excellent; starship shipyard; overhaul; refined fuel",
	'B': "Good; spacecraft shipyard; overhaul; refined fuel",
	'C': "Routine; major-damage repair; unrefined fuel",
	'D': "Poor; minor-damage repair; unrefined fuel",
	'E': "Frontier; beacon only; no facilities",
	'X': "None; no starport",
}

// sizeDesc maps a world size digit to its description and diameter (World Size
// table, p. 25). Size = 2D-2, rerolling a 10 as 1D+9, so the reachable range is
// 0–F.
var sizeDesc = map[int]string{
	0:  "Asteroid Belt",
	1:  "1,000 miles (1,600 km)",
	2:  "2,000 miles (3,200 km)",
	3:  "3,000 miles (4,800 km)",
	4:  "4,000 miles (6,400 km)",
	5:  "5,000 miles (8,000 km)",
	6:  "6,000 miles (9,600 km)",
	7:  "7,000 miles (11,200 km)",
	8:  "8,000 miles (12,800 km)",
	9:  "9,000 miles (14,400 km)",
	10: "10,000 miles (16,000 km)",
	11: "11,000 miles (17,600 km)",
	12: "12,000 miles (19,200 km)",
	13: "13,000 miles (20,800 km)",
	14: "14,000 miles (22,400 km)",
	15: "15,000 miles (24,000 km)",
}

// atmoDesc maps an atmosphere digit to its description (Atmosphere table, 0–F).
var atmoDesc = map[int]string{
	0:  "Vacuum",
	1:  "Trace",
	2:  "Very Thin, Tainted",
	3:  "Very Thin",
	4:  "Thin, Tainted",
	5:  "Thin",
	6:  "Standard",
	7:  "Standard, Tainted",
	8:  "Dense",
	9:  "Dense, Tainted",
	10: "Exotic",
	11: "Corrosive",
	12: "Insidious",
	13: "Dense, High",
	14: "Thin, Low",
	15: "Unusual",
}

// hydroDesc maps a hydrographic digit to its description and water percentage
// (Hydrographics table, 0–A).
var hydroDesc = map[int]string{
	0:  "Desert World",
	1:  "10% Water",
	2:  "20% Water",
	3:  "30% Water",
	4:  "40% Water",
	5:  "50% Water",
	6:  "60% Water",
	7:  "70% Water",
	8:  "80% Water",
	9:  "90% Water",
	10: "Water World",
}

// popDesc maps a population digit to its description (Population table). The digit
// is an exponent of 10.
var popDesc = map[int]string{
	0:  "Unpopulated",
	1:  "Tens (10^1)",
	2:  "Hundreds (10^2)",
	3:  "Thousands (10^3)",
	4:  "Ten Thousands (10^4)",
	5:  "Hundred Thousands (10^5)",
	6:  "Millions (10^6)",
	7:  "Ten Millions (10^7)",
	8:  "Hundred Millions (10^8)",
	9:  "Billions (10^9)",
	10: "Ten Billions (10^10)",
	11: "Hundred Billions (10^11)",
	12: "Trillions (10^12)",
	13: "Ten Trillions (10^13)",
	14: "Hundred Trillions (10^14)",
	15: "Quadrillions (10^15)",
}

// govDesc maps a government digit to its description (Government table, 0–F).
var govDesc = map[int]string{
	0:  "No Government Structure",
	1:  "Company / Corporation",
	2:  "Participating Democracy",
	3:  "Self-Perpetuating Oligarchy",
	4:  "Representative Democracy",
	5:  "Feudal Technocracy",
	6:  "Captive Government / Colony",
	7:  "Balkanization",
	8:  "Civil Service Bureaucracy",
	9:  "Impersonal Bureaucracy",
	10: "Charismatic Dictatorship",
	11: "Non-Charismatic Dictatorship",
	12: "Charismatic Oligarchy",
	13: "Religious Dictatorship",
	14: "Religious Autocracy",
	15: "Totalitarian Oligarchy",
}

// lawDesc maps a law level to its description (Law Level table, 0–J). Law =
// Flux + Gov, capped at J (18).
var lawDesc = map[int]string{
	0:  "No Law. No prohibitions",
	1:  "Low Law. WMD and Psi weapons prohibited",
	2:  "Low Law. Portable weapons prohibited",
	3:  "Low Law. Acid, fire, and gas prohibited",
	4:  "Moderate Law. Laser and beam weapons prohibited",
	5:  "Moderate Law. Shock, EMP, rad, mag, grav weapons prohibited",
	6:  "Moderate Law. Machine guns prohibited",
	7:  "Moderate Law. Pistols prohibited",
	8:  "High Law. Open display of weapons prohibited",
	9:  "High Law. Weapons outside the home prohibited",
	10: "Extreme Law. All weapons prohibited",
	11: "Extreme Law. Continental passports required",
	12: "Extreme Law. Unrestricted invasion of privacy",
	13: "Extreme Law. Paramilitary law enforcement",
	14: "Extreme Law. Full-fledged police state",
	15: "Extreme Law. Daily life rigidly controlled",
	16: "Extreme Law. Disproportionate punishment",
	17: "Extreme Law. Legalized oppressive practices",
	18: "Extreme Law. Routine oppression",
}

// techDesc maps a tech level to its era-band description. Book 3 p. 25 gives only
// the TL DM formula, so these edition-invariant Traveller bands are shared with
// the tne package (see the doc). TL = 1D + DMs tops out at 20: the +2 size DM
// (size 0–1) and +2 hydrographics DM (hydro A) can never both apply, because
// size < 2 forces hydrographics to 0. So every reachable TL has a band here.
var techDesc = map[int]string{
	0:  "Pre-Industrial (primitive)",
	1:  "Pre-Industrial (bronze, iron)",
	2:  "Pre-Industrial (printing press)",
	3:  "Pre-Industrial (basic science)",
	4:  "Industrial (internal combustion)",
	5:  "Industrial (mass production)",
	6:  "Pre-Stellar (nuclear power)",
	7:  "Pre-Stellar (miniaturized electronics)",
	8:  "Pre-Stellar (superconductors)",
	9:  "Early Stellar (fusion power)",
	10: "Early Stellar (jump drive)",
	11: "Average Stellar (large starships)",
	12: "Average Stellar (sophisticated robots)",
	13: "Average Stellar (holographic data storage)",
	14: "High Stellar (anti-grav cities)",
	15: "High Stellar (anagathics)",
	16: "High Stellar (global terraforming)",
	17: "Extreme Stellar",
	18: "Extreme Stellar",
	19: "Extreme Stellar",
	20: "Extreme Stellar (the Ancients)",
}

// starportType maps a 2D roll (2–12) to a mainworld starport code (Book 3 p. 24,
// table 1A), indexed by roll-2 so index 0 is a roll of 2. Identical to Classic.
var starportType = [11]byte{'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'E', 'X'}

// navalTarget and scoutTarget are the per-starport "2D <= target" thresholds for
// a base (Book 3 p. 24 table 1A / p. 28 table B). A missing key is target 0,
// meaning that starport can never hold the base and no die is rolled.
var (
	navalTarget = map[byte]int{'A': 6, 'B': 5}
	scoutTarget = map[byte]int{'A': 4, 'B': 5, 'C': 6, 'D': 7}
)

// Tech-level die modifiers (Book 3 p. 25 TL matrix). Summed with 1D and floored
// at 0. Each map is keyed by a characteristic's value; a missing key contributes
// 0. popTechDM is a function because the source's "Pop A+" row is an open range
// (Pop >= 10 all score +4).
var (
	starportTechDM = map[byte]int{'A': 6, 'B': 4, 'C': 2, 'X': -4}
	sizeTechDM     = map[int]int{0: 2, 1: 2, 2: 1, 3: 1, 4: 1}
	atmoTechDM     = map[int]int{0: 1, 1: 1, 2: 1, 3: 1, 10: 1, 11: 1, 12: 1, 13: 1, 14: 1, 15: 1}
	hydroTechDM    = map[int]int{9: 1, 10: 2}
	govTechDM      = map[int]int{0: 1, 5: 1, 13: -2}
)

// popTechDM returns the population tech DM: +1 for Pop 1–5, +2 for Pop 9, +4 for
// Pop A+ (>= 10), else 0.
func popTechDM(pop int) int {
	switch {
	case pop >= 10:
		return 4
	case pop == 9:
		return 2
	case pop >= 1 && pop <= 5:
		return 1
	default:
		return 0
	}
}
