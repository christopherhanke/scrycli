package main

import (
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
	expected := fmt.Errorf("no query arguments given")
	client := &http.Client{}
	_, err := search(client, []string{})
	if err.Error() != expected.Error() {
		t.Fatalf("Empty Search Error failed\nExpected: %v\nActual: %v", expected, err)
	}
}

func TestSearchQuery(t *testing.T) {
	tests := map[string]struct {
		input    []string
		expected string
	}{
		"Test query only name":      {[]string{"sheoldred", "the"}, "sheoldred+the"},
		"Test query with color":     {[]string{"urza", "c:u"}, "urza+c%3Au"},
		"Test query with cmc":       {[]string{"urza", "cmc=4"}, "urza+cmc%3D4"},
		"Test query with quotation": {[]string{"urza", "o:'draw'"}, "urza+o%3A%27draw%27"},
		"Test query with quotes":    {[]string{"urza", "o:\"draw\""}, "urza+o%3A%22draw%22"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := searchQuery(test.input)
			if err != nil {
				t.Fatalf("Query conversion failed: %v", err)
			}
			if !reflect.DeepEqual(output, test.expected) {
				t.Fatalf("Querry result differ\nExpected: %s\nActual: %s", test.expected, output)
			}
		})
	}
}

func TestSearchQueryEmpty(t *testing.T) {
	expected := fmt.Errorf("no query arguments given")
	_, err := searchQuery([]string{})
	if err.Error() != expected.Error() {
		t.Fatalf("Empty Query Error failed\nExpected: %v\nActual: %v", expected, err)
	}

}
