package utils

import (
	"strings"
)

// maskString masks the string based on the given parameters.
// `s` is the input string, `maskChar` is the character used for masking,
// `option` determines whether the start ("first") or the end ("last") of the string is unmasked,
// `unmaskedCount` is the number of characters that will not be masked at the start or end, in the case of a fixed mask, it is the number of characters that will be masked.
func MaskString(s, maskChar, option string, unmaskedCount int) string {
	fixedMask := unmaskedCount
	n := len(s)
	if unmaskedCount > n {
		unmaskedCount = n
	}

	switch option {
	case "first":
		unmasked := s[:unmaskedCount]
		masked := strings.Repeat(maskChar, n-unmaskedCount)
		return unmasked + masked
	case "last":
		unmasked := s[n-unmaskedCount:]
		masked := strings.Repeat(maskChar, n-unmaskedCount)
		return masked + unmasked
	case "fixed":
		masked := strings.Repeat(maskChar, fixedMask)
		return masked
	default:
		return "XXERRORXX" // In case of an invalid option, return the string unchanged.
	}
}

// MaskStringInMap takes a map of strings and applies the maskString function to each value in the map.
// It returns a new map with the masked strings.
func MaskStringInMap(stringsMap map[string]string, maskChar, option string, unmaskedCount int) map[string]string {
	maskedMap := make(map[string]string)
	for key, value := range stringsMap {
		maskedMap[key] = MaskString(value, maskChar, option, unmaskedCount)
	}
	return maskedMap
}

// Function to create a map from a slice for quicker lookup
func CreateMapFromSlice(list []string) map[string]bool {
	strMap := make(map[string]bool)
	for _, item := range list {
		strMap[item] = true
	}
	return strMap
}

// Function to check if a string is in the generated map (and therefore in the original list)
func StringInMap(str string, strMap map[string]bool) bool {
	return strMap[str]
}
