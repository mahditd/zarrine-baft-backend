package utils

import (
	"errors"
	"strings"
)

// NormalizePhone normalizes Iran-focused phone numbers per SRS 3.2.
// Accepted inputs (same stored value):
//   09121234567, +989121234567, 00989121234567,
//   989121234567, 9121234567, plus spaces/dashes/parens.
// Returns 11-digit 09xxxxxxxxx or error.
func NormalizePhone(phone string) (string, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return "", errors.New("invalid phone number")
	}

	// Convert Persian/Arabic digits to English.
	var sb strings.Builder
	sb.Grow(len(phone))
	for _, r := range phone {
		switch {
		case r >= '0' && r <= '9':
			sb.WriteRune(r)
		case r >= '\u06F0' && r <= '\u06F9':
			sb.WriteRune('0' + (r - '\u06F0'))
		case r >= '\u0660' && r <= '\u0669':
			sb.WriteRune('0' + (r - '\u0660'))
		case r == '+':
			sb.WriteRune(r)
		case r == ' ' || r == '-' || r == '(' || r == ')' || r == '/' || r == '.':
			continue
		default:
			return "", errors.New("invalid phone number")
		}
	}
	cleaned := sb.String()

	var national string
	switch {
	case strings.HasPrefix(cleaned, "+98"):
		national = "0" + cleaned[3:]
	case strings.HasPrefix(cleaned, "0098"):
		national = "0" + cleaned[4:]
	case strings.HasPrefix(cleaned, "98"):
		national = "0" + cleaned[2:]
	case strings.HasPrefix(cleaned, "0"):
		national = cleaned
	case len(cleaned) == 10:
		// Mobile without 0 (912...) or landline without 0 (218...).
		national = "0" + cleaned
	default:
		return "", errors.New("invalid phone number")
	}

	// Mobile (09...) or landline (021...). Both are 11 digits starting with 0.
	if len(national) != 11 || !strings.HasPrefix(national, "0") {
		return "", errors.New("invalid phone number")
	}
	for _, r := range national {
		if r < '0' || r > '9' {
			return "", errors.New("invalid phone number")
		}
	}

	return national, nil
}
