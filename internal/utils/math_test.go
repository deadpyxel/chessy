package utils

import "testing"

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"zero returs zero", 0, 0},
		{"positive 1 returns 1", 1, 1},
		{"negative 1 returns 1", -1, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Abs(test.input)
			if result != test.expected {
				t.Errorf("Abs(%d) = %d; want %d", test.input, result, test.expected)
			}
		})
	}
}
