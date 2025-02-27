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
