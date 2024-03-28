package slicer

import (
	"testing"
)

func TestIsSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "TestSlice",
			input:    []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "TestArray",
			input:    [3]int{1, 2, 3},
			expected: false,
		},
		{
			name:     "TestMap",
			input:    map[string]int{"a": 1, "b": 2},
			expected: false,
		},
		{
			name:     "TestString",
			input:    "Hello",
			expected: false,
		},
		{
			name:     "TestInt",
			input:    42,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsSlice(test.input)
			if result != test.expected {
				t.Errorf("Expected IsSlice(%v) to be %v, but got %v", test.input, test.expected, result)
			}
		})
	}
}
