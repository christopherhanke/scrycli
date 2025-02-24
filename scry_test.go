package main

import (
	"testing"
)

func TestSearch(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected error
	}{
		"Test search - Sheoldred": {"Sheoldred the Apocalypse", nil},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			testFields := cleanInput(test.input)
			_, err := search(testFields)
			if err != nil {
				t.Fatalf("Searching failed: %v", err)
			}
		})
	}
}
