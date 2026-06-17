package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Hello, World!", []string{"hello,", "world!"}},
		{"Hello, WORLD!", []string{"hello,", "world!"}},
		{"hello world", []string{"hello", "world"}},
		{"123", []string{"123"}},
		{"", []string{}},
	}
	for _, test := range tests {
		result := cleanInput(test.input)
		for i, word := range result {
			if word != test.expected[i] {
				t.Errorf("cleanInput(%q) = %v; want %v", test.input, result, test.expected)
				break
			}
		}
	}
}
