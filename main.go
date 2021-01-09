package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const htSitemap = "https://www.hindustantimes.com/sitemap/news.xml"
const MaxArticleCap = 186

// flag variables
var paraSelector = ".storyDetails > .detail > p"
var ifStats bool
var ifHelp bool

// struct article information
type article struct {
	Title string
	URL   string
	Para  string
}

// Text return a concatenated string of article title and paragraph
func (a *article) Text() string {
	if (a).Para == "" {
		return fmt.Sprintf("%s\nLoading", a.Title)
	}

	return fmt.Sprintf("%s\n%s", a.Title, a.Para)
}

// array of articles
// this is a fixed sized array
var articles [MaxArticleCap]article

// no. of articles in the array
var articleL int

// struct to store article text and index
// the program uses go routine to fetch data
// this struct will will be the data received from the channel
// and then added to the article array at given index
type artText struct {
	I    int
	Text string
}

// check if focused on paragraph
var isPara bool = false

// exit the program if the err is not nil
func errHandle(msg string, err error) {
	if err != nil {
		fmt.Printf("%v: %v", msg, err)
		os.Exit(0)
	}
}

// get data from sitemap and return formated data
func getSitemapData() {
	// get the document of the sitemap
	doc := fetchDoc(htSitemap)

	a := new(article)

	// query the url element
	doc.Find("url").Each(func(i int, u *goquery.Selection) {
		// get the title from the XML document
		t, _ := u.Find(`news\:news news\:title`).Html()
		// the title is in cdata block
		// <![CDATA[*title*]>
		a.Title = t[11 : len(t)-5]
		a.URL = u.Find("loc").Text()

		// add the article to the array if the total no. of articles is lower than
		// MaxArticleCap
		if i < MaxArticleCap {
			// add the article pointer to the array at the index i
			articles[i] = *a
			// set the article length
			articleL = i
		}
	})
}

// get the data from the URL and return it through the given channel
func scrapeData(url string, i int, ch chan<- artText) {
	at := new(artText)
	// set the index of artText for proper synchronization of the data
	at.I = i
	// check blank URLs
	if url == "" {
		return
	}
	// get document for the URL
	doc := fetchDoc(url)

	var text string

	// query the article element
	q := doc.Find(paraSelector)

	if q.Length() <= 0 {
		// create an error msg
		at.Text = fmt.Sprintf("error: %v: 0 result from \"%v \"(selector) ", url, paraSelector)
		// send the error msg
		ch <- *at
		return
	}

	q.Each(func(i int, t *goquery.Selection) {
		text = strings.TrimSpace(t.Text())
		// add the text to the final text
		// if the text is not blank
		if text != "" {
			// add extra space if the no the 1st element
			if i != 0 {
				at.Text += ("\n\n" + text)
			} else {
				at.Text += text
			}
		}
	})

	// send data back to the channel
	ch <- *at
}

// fetch all the articles in the articles array
func fetchAllArticles(ch chan<- artText) {
	// loop over the articles jumping 16 items
	for i := 0; i < articleL; i++ {
		// go routine to fetch data
		go scrapeData(articles[i].URL, i, ch)

		// sleep for 2 seconds to avoid rejection from the website
		if i%16 == 0 {
			time.Sleep(2 * time.Second)
		}
	}
}

// sync data of all the articles in the article array
func syncArticles(ch <-chan artText) {
	// wait for the data
	for at := range ch {
		// add the data at the proper index
		articles[at.I].Para = at.Text
	}
}

func init() {
	// fetch sitemap data
	getSitemapData()
	articleL++

	// create a channel to pass the article paragraph data
	ch := make(chan artText, 16)

	// go routine to fetch all the article data
	go fetchAllArticles(ch)

	// go routine to received all the data from the request go routine
	go syncArticles(ch)
}

func main() {
	// initialize the ui
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// get the dimensions
	w, h := ui.TerminalDimensions()
	hw := w/2 - 2
	//	hh := h/2 - 2

	if hw > 80 {
		hw = 80
	}

	// array of list for article titles
	articlesLS := make([]string, articleL)
	// add titles to the array
	for i, a := range articles {
		// check if empty
		if a.Title == "" {
			break
		}
		// concatenate the index and the title
		articlesLS[i] = fmt.Sprintf("%03d  %s", i, a.Title)
	}

	// create new widgets
	l := widgets.NewList()
	p := widgets.NewParagraph()

	// set information about the list
	l.Title = " News Articles List "
	l.Rows = articlesLS // assigning the data
	// setting the style of the widget
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.BorderStyle = ui.NewStyle(ui.ColorGreen)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)

	l.WrapText = false
	// setting the size of the widget
	l.SetRect(0, 0, hw, h)

	// set information about the paragraph
	p.Title = " Articles "
	p.Text = articles[0].Text()
	p.WrapText = true
	// setting the size of the widget
	p.SetRect(hw, 0, w, h)

	// render the UI
	ui.Render(l, p)

	previousKey := ""
	uiEvents := ui.PollEvents()

	// loop to trigger input functions
	for {
		e := <-uiEvents // get inputs

		switch e.ID {
		// quit the application
		case "q", "<C-c>":
			return
		// movements
		case "j", "<Down>":
			if !isPara {
				l.ScrollDown()
				p.Text = articles[l.SelectedRow].Text()
			}
		case "k", "<Up>":
			if !isPara {
				l.ScrollUp()
				p.Text = articles[l.SelectedRow].Text()
			}
		case "l", "<Right>":
			// check if para is not selected
			// change the styles of list and para widget
			if !isPara {
				isPara = true

				// change the list style to deselected
				l.TextStyle = ui.NewStyle(ui.ColorWhite)
				l.BorderStyle = ui.NewStyle(ui.ColorWhite)

				// change the para style to selected
				p.TextStyle = ui.NewStyle(ui.ColorYellow)
				p.BorderStyle = ui.NewStyle(ui.ColorGreen)
			}
		case "h", "<Left>":
			// check if para is selected
			// change the styles of list and para widget
			if isPara {
				isPara = false

				// change the list style to selected
				l.TextStyle = ui.NewStyle(ui.ColorYellow)
				l.BorderStyle = ui.NewStyle(ui.ColorGreen)

				// change the para style to deselected
				p.TextStyle = ui.NewStyle(ui.ColorWhite)
				p.BorderStyle = ui.NewStyle(ui.ColorWhite)
			}
		case "c":
			// copy command
			cp := exec.Command("xclip", "-selection", "c")

			// get the article text
			t := articles[l.SelectedRow].Text()
			t = strings.ReplaceAll(t, "\n\n", " ") // replace the new lines

			// use article as an input for copy command
			cp.Stdin = strings.NewReader(t)

			// run the command
			e := cp.Start()
			errHandle("copy-start", e)
			e = cp.Wait()
			errHandle("copy-wait", e)
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		// re-render the interface
		ui.Render(l, p)
	}
}
