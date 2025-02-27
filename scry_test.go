package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestSearch(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []string
	}{
		"Test search - Sheoldred": {"Sheoldred the Apocalypse", []string{"Sheoldred, the Apocalypse"}},
		"Test search - Lotleth":   {"Lotleth", []string{"Lotleth Giant", "Lotleth Troll"}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &http.Client{}
			testFields := cleanInput(test.input)
			output, err := search(client, testFields)
			if err != nil {
				t.Fatalf("Searching failed: %v", err)
			}
			if !reflect.DeepEqual(test.expected, output) {
				t.Fatalf("Search results differ\nExpected: %s\nActual: %s", test.expected, output)
			}
			time.Sleep(time.Millisecond * 100)
		})
	}
}

func TestEmptySearch(t *testing.T) {
	expected := fmt.Errorf("no args given")
	client := &http.Client{}
	_, err := search(client, []string{})
	if errors.Is(err, expected) {
		t.Fatalf("Empty Search Error failed\nExpected: %v\nActual: %v", expected, err)
	}
}
