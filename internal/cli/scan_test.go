package cli

import (
	"testing"
)

func TestParsePorts(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
		hasError bool
	}{
		{"80,443", []int{80, 443}, false},
		{"100-103", []int{100, 101, 102, 103}, false},
		{"invalid", nil, true},
	}

	for _, tt := range tests {
		ports, err := parsePorts(tt.input)
		if tt.hasError && err == nil {
			t.Errorf("input %q: expected error, got none", tt.input)
		}
		if !tt.hasError && err != nil {
			t.Errorf("input %q: unexpected error: %v", tt.input, err)
		}
		if len(ports) != len(tt.expected) {
			t.Errorf("input %q: expected %v, got %v", tt.input, tt.expected, ports)
		}
	}
}
