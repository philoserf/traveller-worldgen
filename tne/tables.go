package tne

// All tables in this file are transcribed from docs/tne/world-generation.md,
// which was read from the Traveller: The New Era core rulebook (GDW, 1993, the
// Imperial UWP Generation flowchart on printed p.186 and the World Building
// Charts on pp.188-190). TNE's Basic Mainworld Generation is nearly identical to
// MegaTraveller's; the differences are noted where they occur.

// beyondRange is the description returned for a characteristic value with no
// table entry. TNE's ranges are complete for every value the dice can produce,
// so this is a safety net that should not surface in practice.
const beyondRange = "(beyond described range)"

// starportDesc maps a starport code to its description (Starport Quality table).
var starportDesc = map[byte]string{
	'A': "Excellent; starship shipyard; overhaul; refined fuel",
	'B': "Good; spacecraft shipyard; overhaul; refined fuel",
	'C': "Routine; major-damage repair; unrefined fuel",
	'D': "Poor; minor-damage repair; unrefined fuel",
	'E': "Frontier; no facilities",
	'X': "None; no starport",
}

// sizeDesc maps a world size digit to its description and diameter (World Size
// table). The 2D-2 roll reaches at most A (10).
var sizeDesc = map[int]string{
	0:  "Asteroid/Planetoid Belt",
	1:  "Small (1,600 km)",
	2:  "Small (3,200 km, Luna)",
	3:  "Small (4,800 km, Mercury)",
	4:  "Small (6,400 km, Mars)",
	5:  "Medium (8,000 km)",
	6:  "Medium (9,600 km)",
	7:  "Medium (11,200 km)",
	8:  "Large (12,800 km, Terra)",
	9:  "Large (14,400 km)",
	10: "Large (16,000 km)",
}

// atmoDesc maps an atmosphere digit to its description (World Atmosphere table,
// 0-F). TNE labels D/E/F as Exotic variants.
var atmoDesc = map[int]string{
	0:  "Vacuum",
	1:  "Vacuum (trace)",
	2:  "Vacuum (very thin, tainted)",
	3:  "Vacuum (very thin)",
	4:  "Thin (tainted)",
	5:  "Thin",
	6:  "Standard",
	7:  "Standard (tainted)",
	8:  "Dense",
	9:  "Dense (tainted)",
	10: "Exotic",
	11: "Exotic (corrosive)",
	12: "Exotic (insidious)",
	13: "Exotic (dense, high)",
	14: "Exotic (thin, low)",
	15: "Exotic (ellipsoid)",
}

// hydroDesc maps a hydrographic digit to its description and water percentage
// (World Hydrographics, 0-A).
var hydroDesc = map[int]string{
	0:  "Desert World (0%)",
	1:  "Dry World (10%)",
	2:  "Dry World (20%)",
	3:  "Wet World (30%)",
	4:  "Wet World (40%)",
	5:  "Wet World (50%)",
	6:  "Wet World (60%)",
	7:  "Wet World (70%)",
	8:  "Wet World (80%)",
	9:  "Wet World (90%)",
	10: "Water World (100%)",
}

// popDesc maps a population digit to its description (World Population table). The
// digit is an exponent of 10. TNE labels the smallest populations
// "Inconsequential".
var popDesc = map[int]string{
	0:  "Inconsequential (< ten)",
	1:  "Inconsequential (tens)",
	2:  "Inconsequential (hundreds)",
	3:  "Low (thousands)",
	4:  "Low (ten-thousands)",
	5:  "Low (hundred-thousands)",
	6:  "Moderate (millions)",
	7:  "Moderate (ten-millions)",
	8:  "Moderate (hundred-millions)",
	9:  "High (billions)",
	10: "High (ten-billions)",
}

// govDesc maps a government digit to its description (World Government table,
// 0-F).
var govDesc = map[int]string{
	0:  "No Government Structure",
	1:  "Company/Corporation",
	2:  "Participating Democracy",
	3:  "Self-Perpetuating Oligarchy",
	4:  "Representative Democracy",
	5:  "Feudal Technocracy",
	6:  "Captive Government / Colony",
	7:  "Balkanization",
	8:  "Civil Service Bureaucracy",
	9:  "Impersonal Bureaucracy",
	10: "Charismatic Dictator",
	11: "Non-Charismatic Leader",
	12: "Charismatic Oligarchy",
	13: "Religious Dictatorship",
	14: "Religious Autocracy",
	15: "Totalitarian Oligarchy",
}

// lawDesc maps a law level to its description (World Law Level table, 0-L).
var lawDesc = map[int]string{
	0:  "No law (no prohibitions)",
	1:  "Low law (body pistols, explosives, poison gas prohibited)",
	2:  "Low law (portable energy weapons prohibited)",
	3:  "Low law (machine guns, automatic rifles prohibited)",
	4:  "Moderate law (light assault weapons prohibited)",
	5:  "Moderate law (personal concealable weapons prohibited)",
	6:  "Moderate law (all firearms except shotguns prohibited)",
	7:  "Moderate law (shotguns prohibited)",
	8:  "High law (blade weapons controlled, no open display)",
	9:  "High law (weapon possession outside the home prohibited)",
	10: "Extreme law (weapon possession prohibited)",
	11: "Extreme law (rigid control of civilian movement)",
	12: "Extreme law (unrestricted invasion of privacy)",
	13: "Extreme law (paramilitary law enforcement)",
	14: "Extreme law (full-fledged police state)",
	15: "Extreme law (all facets of daily life rigidly controlled)",
	16: "Extreme law (severe punishment for petty infractions)",
	17: "Extreme law (legalized oppressive practices)",
	18: "Extreme law (routinely oppressive and restrictive)",
	19: "Extreme law (excessively oppressive and restrictive)",
	20: "Extreme law (totally oppressive and restrictive)",
}

// techDesc maps a tech level to its description (Technology Level table, 0-L).
// TNE places the Imperium's typical level at D and the Ancients at L.
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
	13: "Average Stellar (holographic data storage; the Imperium)",
	14: "High Stellar (anti-grav cities)",
	15: "High Stellar (anagathics)",
	16: "High Stellar (global terraforming)",
	17: "Extreme Stellar",
	18: "Extreme Stellar",
	19: "Extreme Stellar",
	20: "Extreme Stellar (the Ancients)",
}

// starportByNature maps a subsector nature to its 2D starport-result column,
// indexed by (roll - 2) so index 0 is a roll of 2 and index 10 is a roll of 12
// (Starport table). Identical to MegaTraveller's.
var starportByNature = map[Nature][11]byte{
	Backwater: {'A', 'A', 'B', 'B', 'C', 'C', 'C', 'D', 'E', 'E', 'X'},
	Standard:  {'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'E', 'X'},
	Mature:    {'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'E', 'E'},
	Cluster:   {'A', 'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'X'},
}

// baseThrow is the per-starport 2D target for each base kind. A zero target means
// that starport can never have the base (no die is rolled). Identical to
// MegaTraveller. Way Station (code B) is deferred.
type baseThrow struct {
	naval    int
	scout    int
	military int
}

var baseThrows = map[byte]baseThrow{
	'A': {naval: 8, scout: 10, military: 10},
	'B': {naval: 8, scout: 9, military: 9},
	'C': {scout: 8, military: 8},
	'D': {scout: 7},
}

// Technology die modifiers (Technology DM matrix). Each is keyed by a
// characteristic's value; a missing key contributes 0. Summed with 1D and
// floored at 0 to yield the tech level. TNE adds three cells over MegaTraveller:
// Starport F (+1), Government E (-1), and Government F (-1). Starport A is +6.
//
// The Starport 'F' entry mirrors the printed source matrix but is never consulted
// here: a mainworld starport is only ever A-E or X (F/G/H are spaceport types,
// out of scope). Government E/F, by contrast, are reachable (government runs
// 0-F).
var (
	starportTechDM = map[byte]int{'A': 6, 'B': 4, 'C': 2, 'F': 1, 'X': -4}
	sizeTechDM     = map[int]int{0: 2, 1: 2, 2: 1, 3: 1, 4: 1}
	atmoTechDM     = map[int]int{0: 1, 1: 1, 2: 1, 3: 1, 10: 1, 11: 1, 12: 1, 13: 1, 14: 1, 15: 1}
	hydroTechDM    = map[int]int{9: 1, 10: 2}
	popTechDM      = map[int]int{1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 9: 2, 10: 4}
	govTechDM      = map[int]int{0: 1, 5: 1, 13: -2, 14: -1, 15: -1}
)
