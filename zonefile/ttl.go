package zonefile

import "strconv"

// StringToTTL parses things like 2w, 2m, etc, and returns the time in seconds.
func StringToTTL(token string) (uint32, bool) {
	// Try to parse as a plain number first
	if seconds, err := strconv.ParseUint(token, 10, 32); err == nil {
		return uint32(seconds), true
	}
	var s, i uint32
	for _, c := range token {
		switch c {
		case 's', 'S':
			s += i
			i = 0
		case 'm', 'M':
			s += i * 60
			i = 0
		case 'h', 'H':
			s += i * 60 * 60
			i = 0
		case 'd', 'D':
			s += i * 60 * 60 * 24
			i = 0
		case 'w', 'W':
			s += i * 60 * 60 * 24 * 7
			i = 0
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			i *= 10
			i += uint32(c) - '0'
		default:
			return 0, false
		}
	}
	return s + i, true
}
