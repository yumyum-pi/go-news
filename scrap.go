package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

// fetch all the articleList in the articleList array
func fetchAllArticleList() {
	// loop over the articleList jumping 16 items
	for i := 0; i < articleL; i++ {
		// go routine to fetch data
		go func(url string, index int) {
			fetchDoc(url, func(doc *goquery.Document) {
				// check blank URLs
				if url == "" {
					return
				}

				var text string
				at := new(artText)

				// set the index of artText for proper synchronization of the data
				at.I = index

				// query the article element
				q := doc.Find(paraSelector)

				// check the no. of elements
				if q.Length() <= 0 {
					// create an error msg
					at.Text = fmt.Sprintf(
						"error: %v: 0 result from \"%v \"(selector) ",
						url,
						paraSelector)
					// send the error msg
					ch <- *at
					return
				}

				// scrap the data
				q.Each(func(j int, t *goquery.Selection) {
					text = strings.TrimSpace(t.Text())
					// add the text to the final text
					// if the text is not blank
					if text != "" {
						// add extra space if the no the 1st element
						if j != 0 {
							at.Text += ("\n\n" + text)
						} else {
							at.Text += text
						}
					}
				})

				// send data back to the channel
				ch <- *at
			})
		}(articleList[i].URL, i)

		// sleep for 2 seconds to avoid rejection from the website
		if i%16 == 0 {
			time.Sleep(2 * time.Second)
		}
	}
}

// get the list of articles from sitemap
func sitemapScraper(doc *goquery.Document) {
	// constants
	const elemS = "url"                                            // element selector
	const titleS = `news\:news news\:title`                        // title selector
	const pubS = `news\:news news\:publication news\:name`         // title selector
	const pubLangS = `news\:news news\:publication news\:language` // title selector
	const pubDate = `news\:news news\:publication_date`            // title selector
	const urlS = "loc"                                             // URL selector

	a := new(article)

	// check the no. of elements
	q := doc.Find(elemS)

	if q.Length() <= 0 {
		// create an error msg
		errHandle("", fmt.Errorf("error: %v: 0 result from \"%v \"(selector) ", "sitemapScraper", elemS))
		return
	}

	// loop thought each element and extract data
	q.Each(func(i int, u *goquery.Selection) {
		// get the title from the XML document
		t, _ := u.Find(titleS).Html()
		// the title is in cdata block
		// which is not parsed by the, need to manually select the text
		// <![CDATA[*title*]>
		a.Title = t[11 : len(t)-5]
		a.URL = u.Find(urlS).Text()
		a.Date = u.Find(pubDate).Text()
		a.Language = u.Find(pubLangS).Text()
		a.Source = u.Find(pubS).Text()

		// add the article to the array if the total no. of articleList is lower than
		// MaxArticleCap
		if i < MaxArticleCap {
			// add the article pointer to the array at the index i
			articleList[i] = *a
			// set the article length
			articleL = i
		}
	})
}
