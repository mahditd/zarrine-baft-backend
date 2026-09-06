package utils_test

import (
	"testing"

	"github.com/mahditd/zarrine-baft-backend/internal/utils"
)

func TestNormalizePersian(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "  مشكي   ",
			expected: "مشکی",
		},
		{
			input:    "پیراهن   مردانه   آبی  ",
			expected: "پیراهن مردانه آبی",
		},
		{
			input:    "كت   زمستاني",
			expected: "کت زمستانی",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		got := utils.NormalizePersian(tt.input)
		if got != tt.expected {
			t.Errorf("NormalizePersian(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestNormalizeEnglish(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "  Black  ",
			expected: "black",
		},
		{
			input:    "WINTER   COAT",
			expected: "winter coat",
		},
	}

	for _, tt := range tests {
		got := utils.NormalizeEnglish(tt.input)
		if got != tt.expected {
			t.Errorf("NormalizeEnglish(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}
}
