package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	mui "github.com/yumyum-pi/go-news/pkg/ui"
)

const htSitemap = "https://www.hindustantimes.com/sitemap/news.xml"

// flag variables
var paraSelector = ".storyDetails > .detail > p"
var ifStats bool
var ifHelp bool

// list width
const listWidth = 40

// MaxArticleCap defines the maximum capacity of the article array
const MaxArticleCap = 186

// array of articleList
// this is a fixed sized array
var articleList [MaxArticleCap]article

// no. of articleList in the array
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

var ch chan artText

// exit the program if the err is not nil
func errHandle(msg string, err error) {
	if err != nil {
		fmt.Printf("%v: %v", msg, err)
		os.Exit(0)
	}
}

// sync data of all the articleList in the article array
func syncarticleList() {
	// wait for the data
	for at := range ch {
		// add the data at the proper index
		articleList[at.I].Para = at.Text
	}
}

func init() {
	// fetch sitemap data
	// get the document of the sitemap
	fetchDoc(htSitemap, sitemapScraper)
	articleL++

	// create a channel to pass the article paragraph data
	ch = make(chan artText, 16)

	// go routine to fetch all the article data
	go fetchAllArticleList()

	// go routine to received all the data from the request go routine
	go syncarticleList()
}

func getColumnWidth(w int) int {
	return w - listWidth
}

func main() {
	// initialize the ui
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// get the dimensions
	w, h := ui.TerminalDimensions()

	// array of list for article titles
	articleListLS := make([]string, articleL)
	// add titles to the array
	for i, a := range articleList {
		// check if empty
		if a.Title == "" {
			break
		}
		// concatenate the index and the title
		articleListLS[i] = a.Title
	}

	// create new widgets
	l := widgets.NewList()
	p := mui.NewParagraph()

	// set information about the list
	l.Title = "List "
	l.Rows = articleListLS // assigning the data
	// setting the style of the widget
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.BorderStyle = ui.NewStyle(ui.ColorGreen)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)

	l.WrapText = false
	// setting the size of the widget
	l.SetRect(0, 0, listWidth, h)

	// set information about the paragraph
	p.Title = "Article"
	p.SetPara(articleList[0].Text(), 0)
	p.WrapText = true
	// setting the size of the widget
	p.SetRect(listWidth, 0, w, h)

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
				p.SetPara(articleList[l.SelectedRow].Text(), 0)
				p.ScrollPosition = articleList[l.SelectedRow].ScrollPositin
			} else {
				p.ScrollDown()
			}
		case "k", "<Up>":
			if !isPara {
				l.ScrollUp()
				p.SetPara(articleList[l.SelectedRow].Text(), 0)
				p.ScrollPosition = articleList[l.SelectedRow].ScrollPositin
			} else {
				p.ScrollUp()
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
			t := articleList[l.SelectedRow].Text()
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
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			p.SetRect(listWidth, 0, payload.Width, payload.Height)
			l.SetRect(0, 0, listWidth, payload.Height)
			ui.Clear()
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
