package utils

// ExtractUnitNumber extracts the unit number from a unit code
// e.g. "FIT3152" -> "3152"
func ExtractUnitNumber(code string) string {
	var number string
	for _, ch := range code {
		if ch >= '0' && ch <= '9' {
			number += string(ch)
		}
	}
	return number
}
