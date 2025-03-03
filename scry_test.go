package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSearch(t *testing.T) {
	// This test only tests on the name of the cards in result to search (not storing plain card data in test)
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
			var stringOuput []string
			for _, value := range output {
				stringOuput = append(stringOuput, value.Name)
			}
			for i := range len(stringOuput) {
				if stringOuput[i] == test.expected[i] {
					continue
				} else {
					t.Fatalf("Search results differ\nExpected: %s\nActual %s", test.expected[i], stringOuput[i])
				}
			}
		})
		time.Sleep(time.Millisecond * 100)
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

func TestRandomCardColor(t *testing.T) {
	tests := map[string]struct {
		input    []string
		expected string
	}{
		"Test random green": {[]string{"c:g"}, "g"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &http.Client{}
			output, err := randomCard(client, test.input)

			if err != nil {
				t.Fatalf("Searching for random failed: %v", err)
			}
			if !strings.Contains(strings.ToLower(output.ManaCost), test.expected) {
				t.Fatalf("Query result differ\nExpected: %s\nActual: %s", test.expected, output.ManaCost)
			}
		})
		time.Sleep(time.Millisecond * 100)
	}
}
