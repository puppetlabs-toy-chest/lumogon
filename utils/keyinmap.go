package utils

// KeyInMap checks if the supplied string s, is present as a key
// in the supplied map[string]int (which is being used as a set)
func KeyInMap(s string, m map[string]int) bool {
	if _, ok := m[s]; ok {
		return true
	}
	return false
}
