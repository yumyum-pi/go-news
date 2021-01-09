package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// count the no of HTML children
func numberOfChild(n *html.Node) int {
	if n == nil {
		return -1
	}
	count := 0

	// c is the current child
	// loop until c is null and iterate the next sibling
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		count++
	}

	return count
}

// map the character in the given text
func mapChar(s *string) map[string]int {
	m := make(map[string]int)
	for i := range *s {
		m[string((*s)[i])]++
	}
	return m
}

// display stats
func stats(s *string) {
	// get the no. of characters
	c := len(*s)
	m := mapChar(s)

	// divide the string to words
	words := strings.Fields(*s)

	fmt.Printf("no. of \ncharacters: %d\nwords: %d\nmap: %v", c, len(words), m)

	// display in a different format
	// sort the data
	// print the map of the character
	for i := range m {
		fmt.Println(i, ":", m[i])
	}
}
