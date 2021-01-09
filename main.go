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
var paraSelector = ".storyDetail > .details > p"
var ifStats bool
var ifHelp bool

type article struct{
	Title string
	URL string
	Para string
}

type artText struct {
	I int
	Text string
}

func (a *article) Text() string {
	if((a).Para == "") {
		return fmt.Sprintf("%s\nLoading", a.Title)
	}

	return fmt.Sprintf("%s\n%s", a.Title, a.Para)
}

var articles [MaxArticleCap]article;
var articleL int

var isPara bool = false

// write a better error handler
func errHandle(msg string, err error) {
	if err != nil {
		fmt.Printf("%v: %v",msg, err)
		os.Exit(0)
	}
}

// get data from sitemap and return formated data
func getSitemapData() {
	// get the document of the sitemap
	doc := fetchDoc(htSitemap)

	a := new(article)

	doc.Find("url").Each(func(i int, u *goquery.Selection){
		// get the title from the XML document
		t, _ := u.Find(`news\:news news\:title`).Html()
		// the title is in cdata block
		// <![CDATA[*title*]>
		a.Title = t[11:len(t)-5]
		a.URL = u.Find("loc").Text()

		// add the article to the array if the total no. of articles is lower than
		// MaxArticleCap
		if ( i < MaxArticleCap ) {
			// add the article pointer to the array at the index i
			articles[i]=*a
			// set the article length
			articleL = i
		}
	})
}

func main() {
	getSitemapData()
	articleL++

	// create a channel
	ch := make( chan artText, 16)

	go func() {
		// loop over the articles jumping 40 items
		for i:= 0; i < articleL; i ++  {
			go scrapeData(articles[i].URL, i, ch)
			if i % 16 == 0 {
				time.Sleep(2 * time.Second)
			}
		}
	}()

	go func () {
		// wait for the data
		for at := range ch {
			articles[at.I].Para = at.Text
		}
	}()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	w, h:= ui.TerminalDimensions()
	hw := w/2 - 2
//	hh := h/2 - 2

	if (hw > 80){
		hw = 80
	}

	articlesLS := make([]string, articleL)
	for i, a := range articles {
		// check if empty
		if a.Title == "" {
			break
		}
		articlesLS[i] = fmt.Sprintf("%03d  %s", i, a.Title)
	}

	l := widgets.NewList()
	p := widgets.NewParagraph()

	// set information about the list
	l.Title = " News Articles List "
	l.Rows = articlesLS
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.BorderStyle = ui.NewStyle(ui.ColorGreen)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	l.WrapText = false
	l.SetRect(0, 0, hw, h)

	// set information about the paragraph
	p.Title = " Articles "
	p.Text = articles[0].Text()
	p.WrapText = true
	p.SetRect(hw,0,w, h)
	ui.Render(l,p)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
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
			if !isPara {
				isPara = true
				l.TextStyle = ui.NewStyle(ui.ColorWhite)
				l.BorderStyle = ui.NewStyle(ui.ColorWhite)
				p.TextStyle = ui.NewStyle(ui.ColorYellow)
				p.BorderStyle = ui.NewStyle(ui.ColorGreen)
			}
		case "h", "<Left>":
			if isPara {
				isPara = false
				l.TextStyle = ui.NewStyle(ui.ColorYellow)
				l.BorderStyle = ui.NewStyle(ui.ColorGreen)
				p.TextStyle = ui.NewStyle(ui.ColorWhite)
				p.BorderStyle = ui.NewStyle(ui.ColorWhite)
			}
		case "c":
			// copy command
			cp := exec.Command("xclip", "-selection", "c")

			t := articles[l.SelectedRow].Text()

			t = strings.ReplaceAll(t,"\n\n", " ")
			// use article as an input for copy command
			cp.Stdin = strings.NewReader(t)
			e := cp.Start()
			errHandle("copy-start",e)
			e = cp.Wait()
			errHandle("copy-wait",e)
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

		ui.Render(l,p)
	}
}



// get the data from the URL
func scrapeData(url string, i int, ch chan artText){
	at := new(artText)
	at.I = i
	// check blank URLs
	if url == "" {
		return
	}
	// get document for the URL
	doc := fetchDoc(url)

	var text string
	// query the article element
	doc.Find(paraSelector).Each(func(i int, t *goquery.Selection) {
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
