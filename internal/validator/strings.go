package validator

import "unicode"

// Read input string and convert first string letter to Upper
func ToCapital(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}
