package ui

import (
	"image"

	ui "github.com/gizak/termui/v3"
)

// Paragraph struct for creating paragraphs
type Paragraph struct {
	ui.Block
	TextStyle      ui.Style
	WrapText       bool
	ScrollPosition int
	Rows           [][]ui.Cell
	RowsL          int
}

// NewParagraph returns pointer to a Paragraph
func NewParagraph() *Paragraph {
	return &Paragraph{
		Block:     *ui.NewBlock(),
		TextStyle: ui.Theme.Paragraph.Text,
		WrapText:  true,
	}
}

// SetPara sets initial variable for the paragraph
func (para *Paragraph) SetPara(text string, scrollPosition int) {
	// get scrollPosition
	para.ScrollPosition = scrollPosition
	// convert paragraph into cell with styles
	cells := ui.ParseStyles(text, para.TextStyle)
	// warp the text
	if para.WrapText {
		cells = ui.WrapCells(cells, uint(para.Inner.Dx()))
	}

	// split the text into rows
	para.Rows = ui.SplitCells(cells, '\n')
	para.RowsL = len(para.Rows)
}

// Draw renders the Paragraph
func (para *Paragraph) Draw(buf *ui.Buffer) {
	// Block manages size, position, border, and title.
	para.Block.Draw(buf)

	if para.Rows == nil {
		return
	}

	l := len(para.Rows)
	// render the paragraph row by row
	for y := 0; y < l; y++ {

		// break if the scroll position exceeds the paragraph length
		if y+para.Inner.Min.Y >= para.Inner.Max.Y {
			break
		}

		row := para.Rows[y+para.ScrollPosition]
		//render the line
		row = ui.TrimCells(row, para.Inner.Dx())
		for _, cx := range ui.BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(para.Inner.Min))
		}
	}
}

// ScrollAmount charge scroll position relative to the given amount
func (para *Paragraph) ScrollAmount(a int) {
	// calculate the scroll position
	toScroll := a + para.ScrollPosition
	// calculate the lowest allowed scrabble position
	bottomScroll := para.RowsL - para.Inner.Dy()

	// should not be 0
	if bottomScroll < 0 {
		// there is no room for scrolling
		return
	}

	// making sure that the toScroll is within the boundaries
	if toScroll < 0 {
		toScroll = 0
	} else if toScroll > bottomScroll {
		toScroll = bottomScroll
	}

	// update the scroll position
	para.ScrollPosition = toScroll
}

// ScrollUp will scroll up the paragraph
func (para *Paragraph) ScrollUp() {
	para.ScrollAmount(-1)
}

// ScrollDown will scroll down the paragraph
func (para *Paragraph) ScrollDown() {
	para.ScrollAmount(+1)
}
