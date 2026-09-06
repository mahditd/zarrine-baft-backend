package utils

import (
	"regexp"
	"strings"
)

var extraSpacesRegex = regexp.MustCompile(`\s+`)

// NormalizePersian normalizes Persian text according to SRS Section 4.6 & 8.1:
// - ي (\u064A) -> ی (\u06CC)
// - ك (\u0643) -> ک (\u06A9)
// - Remove leading/trailing spaces
// - Collapse multiple spaces into a single space
func NormalizePersian(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	// Arabic Yeh to Persian Yeh
	text = strings.ReplaceAll(text, "\u064A", "\u06CC")
	// Arabic Yeh with Hamza to Persian Yeh with Hamza or Yeh
	text = strings.ReplaceAll(text, "\u0649", "\u06CC") // Alef Maksura to Persian Yeh

	// Arabic Kaf to Persian Keheh
	text = strings.ReplaceAll(text, "\u0643", "\u06A9")

	// Collapse whitespace
	text = extraSpacesRegex.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

// NormalizeEnglish normalizes English text:
// - Trim spaces
// - Collapse whitespace
// - Case-insensitive lowercasing (SRS 4.6)
func NormalizeEnglish(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	text = extraSpacesRegex.ReplaceAllString(text, " ")
	return strings.ToLower(text)
}
