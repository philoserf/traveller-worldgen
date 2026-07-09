// Package ehex implements Traveller extended-hexadecimal digits: 0-9 then the
// uppercase letters, omitting I and O to avoid confusion with 1 and 0. Every
// Traveller edition uses this notation for UWP characteristics, so it lives in
// its own edition-independent package.
package ehex

// Digits is the extended-hex alphabet.
const Digits = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"

// Encode returns the extended-hex digit for v, or '?' when v falls outside the
// alphabet's range (0..len(Digits)-1).
func Encode(v int) byte {
	if v < 0 || v >= len(Digits) {
		return '?'
	}
	return Digits[v]
}

// Decode returns the integer value of an extended-hex digit and whether it was
// recognized — the inverse of Encode for valid uppercase input.
func Decode(c byte) (int, bool) {
	for i := range len(Digits) {
		if Digits[i] == c {
			return i, true
		}
	}
	return 0, false
}
