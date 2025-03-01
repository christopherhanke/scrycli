package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const apiURL = "https://api.scryfall.com"

func search(client *http.Client, args []string) ([]string, error) {
	// build API string and parse to URL
	queryString, err := searchQuery(args)
	if err != nil {
		return nil, err
	}
	searchString := apiURL + "/cards/search?q=" + queryString

	searchURL, err := url.Parse(searchString)
	if err != nil {
		return nil, err
	}

	// call API and close response
	req, err := http.NewRequest("GET", searchURL.String(), nil)
	if err != nil {
		return nil, err
	}
	setRequest(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check response status code for too many requests
	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(time.Millisecond * 100)
		return nil, fmt.Errorf("too many requests")
	}

	// unmarshal json data from response and list the names of the cards
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var scryresp scryResponse
	err = json.Unmarshal(data, &scryresp)
	if err != nil {
		return nil, err
	}
	var returnlist []string
	for _, entry := range scryresp.Data {
		returnlist = append(returnlist, entry.Name)
	}

	return returnlist, nil
}

// build search query string from string list
func searchQuery(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("no query arguments given")
	}
	name := args[0]
	if len(args) > 1 {
		for _, arg := range args[1:] {
			name += " " + arg
		}
	}
	return url.QueryEscape(name), nil
}

func searchHelp() {
	fmt.Println("here comes help")
}
