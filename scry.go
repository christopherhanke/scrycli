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

func search(client *http.Client, args []string) ([]card, error) {
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
	if err = checkResponseCode(resp); err != nil {
		return nil, err
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

	var returnlist []card
	for _, entry := range scryresp.Data {
		returnlist = append(returnlist, entry)
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
	fmt.Print("usage: search <search term>\n\n")
	fmt.Println("You can just write a search term like a card name and/or you can use queries like color and manacost.")
	fmt.Print("Search queries work with a query key and a colon (e.g. 'c:r' for color red) or comparison like '<='/'>=' etc.\n\n")
	fmt.Println("These are some of the queries that are possibel:")
	fmt.Println(" color: [c:wu] for cards that are white and blue. [id:wu] checks for color identity white and blue. You can use 'w, u, b, r, g'.")
	fmt.Println(" type: [t:creature] for creature cards. Any supertype, card type or subtype is possible.")
	fmt.Println(" oracle: [o:'text'] searches fot <text> in oracle text.")
	fmt.Println(" mana: [m:g] for cards with one green mana in their mana costs. You can staple these like [c:gguu] for two green and two blue mana pips. You can use 'w, u, b, r, g'.")
	fmt.Println(" manavalue: [mv<=4] for cards with mana value lower or equal to four.")
	fmt.Println(" rarity: [r:c] for common cards. You can use compare operands and 'c, u, r, m'.")
	fmt.Println(" set: [set:khm] for cards from Kaldheim. Use set with their magic set code.")
	fmt.Println(" format: [f:pauper] for cards legal in pauper. Supported: standard, pioneer, modern, legacy, pauper, vintage, commander.")
}

func randomCard(client *http.Client) (card, error) {
	// build API string
	urlString := apiURL + "/cards/random"

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return card{}, err
	}
	setRequest(req)
	resp, err := client.Do(req)
	if err != nil {
		return card{}, err
	}
	defer resp.Body.Close()

	// check response status code for too many requests
	if err = checkResponseCode(resp); err != nil {
		return card{}, err
	}

	//unmarshal json data from response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return card{}, err
	}
	var result card
	err = json.Unmarshal(data, &result)
	if err != nil {
		return card{}, err
	}
	return result, nil
}

// check HTTP Response code for too many requests
func checkResponseCode(resp *http.Response) error {
	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(time.Millisecond * 100)
		return fmt.Errorf("too many requests")
	}
	return nil
}
