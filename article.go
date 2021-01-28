package main

import "fmt"

// struct article information
type article struct {
	Title         string
	Date          string
	Source        string
	Language      string
	URL           string
	Para          string
	ScrollPositin int
}

// Text return a concatenated string of article title and paragraph
func (a *article) Text() string {
	// info about the article
	t := fmt.Sprintf(
		"%s\nTime:%s    Date:%s    Source:%s    Language:%s",
		a.Title,
		a.Date[11:19], // split the date into time
		a.Date[0:10],  // split the date into date
		a.Source,
		a.Language,
	)
	// loading screen
	if (a).Para == "" {
		return fmt.Sprintf("%s\n\nLoading...", t)
	}

	return fmt.Sprintf("%s\n\n%s", t, a.Para)
}
