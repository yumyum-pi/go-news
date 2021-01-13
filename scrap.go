package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// function identifier for scraping data
type scraper func(*goquery.Document)

func fetchDoc(url string, scrap scraper) {
	// get the response from the URL
	res, err := http.Get(url)
	errHandle("scrap: http get error", err)

	// handle unsuccessful status code
	if res.StatusCode != 200 {
		errHandle("scrap: error code from the server", fmt.Errorf(
			"%v:%v",
			res.StatusCode,
			htSitemap,
		))
	}

	defer res.Body.Close()

	// create a new document from the response
	doc, err := goquery.NewDocumentFromReader(res.Body)
	errHandle("scrap: new document from reader", err)

	// run the scraper function
	scrap(doc)
}
