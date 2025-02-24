package main

import (
	"fmt"
	"net/http"
	"net/url"
)

const apiURL = "https://api.scryfall.com"

func search(args []string) ([]string, error) {
	// build API string, escape query and parse to URL
	searchString := apiURL + "/cards/search?q="
	if len(args) < 1 {
		return nil, fmt.Errorf("no args given")
	}
	name := args[0]
	if len(args) > 1 {
		for _, arg := range args[1:] {
			name += " " + arg
		}
	}
	searchString += url.QueryEscape(name)

	searchURL, err := url.Parse(searchString)
	if err != nil {
		return nil, err
	}

	// call API and close response
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ScryCLI/0.1")
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("URL: %s\n", searchURL.String())
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Header: %s\n", resp.Header.Get("content-type"))

	return nil, nil
}
