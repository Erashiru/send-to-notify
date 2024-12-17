package managers

import (
	"testing"
)

func TestRemoveNonAlphaNumericSymbols(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Chicken Republic Ritz Palace", "Chicken Republic Ritz Palace"},
		{"Chiko | Алматы Абылайхана", "Chiko Алматы Абылайхана"},
		{"Del Cappuccino (Гагарина)", "Del Cappuccino Гагарина"},
		{"123 456", "123 456"},
		{"!@#$%^", ""},
		{"", ""},
		{"abc123", "abc123"},
	}

	for _, test := range tests {
		result := reduceSpaces(removeNonAlphaNumericSymbols(test.input))
		if result != test.expected {
			t.Errorf("For input %q, expected %q, but got %q", test.input, test.expected, result)
		}
	}
}
