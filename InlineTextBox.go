package stagui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Like a text box, but only has one
// line
type InlineTextBox struct {
	X, Y, W, H float64

	Text      string
	TextColor color.Color

	BgColor color.Color

	// Is the user attempting to type
	// inside the textbox
	Active bool

	// Where in the textbox the user is
	// typing
	KeyPos int
}

func (itb *InlineTextBox) Draw(screen *ebiten.Image, fh FontHandler) {
	vector.DrawFilledRect(screen, float32(itb.X), float32(itb.Y), float32(itb.W), float32(itb.H), itb.BgColor, false)

	fh.DrawText(screen, itb.Text, 30, itb.X+4, itb.Y+2, fh.GetFont("textBox"), &text.DrawOptions{})

	// Draw a line at the bottom of the textbox
	if itb.Active {
		vector.StrokeLine(screen, float32(itb.X), float32(itb.Y+itb.H), float32(itb.X+itb.W), float32(itb.Y+itb.H), 1, color.White, false)
	}
}

func (itb *InlineTextBox) CheckClick(mx, my int) bool {
	x := float64(mx)
	y := float64(my)

	clicked := itb.X <= x && x <= itb.X+itb.W &&
		itb.Y <= y && y <= itb.Y+itb.H

	// Effectively an xor operation.
	// If we're active and clicked, set
	// active to false
	itb.Active = clicked != itb.Active

	return clicked
}

func (itb *InlineTextBox) Update() {
	if !itb.Active {
		return
	}

	keyText, key := handleKey()
	if keyText == "None" {
		return
	}

	switch key {
	// Clean up rogue inputs
	case ebiten.KeyInsert:
	case ebiten.KeyPageUp:
	case ebiten.KeyPageDown:
	case ebiten.KeyEscape:
	case ebiten.KeyCapsLock:
	case ebiten.KeyControl:
	case ebiten.KeyAlt:
	case ebiten.KeyNumLock:
	case ebiten.KeyContextMenu:

	case ebiten.KeyEnter:
		itb.Active = false
	case ebiten.KeyBackspace:
		if itb.KeyPos == 0 {
			break
		}

		itb.Text = itb.Text[:itb.KeyPos-1] + itb.Text[itb.KeyPos:]
		itb.KeyPos--
	case ebiten.KeyDelete:
		if len(itb.Text) == 0 || itb.KeyPos == len(itb.Text) {
			break
		}

		itb.Text = itb.Text[:itb.KeyPos] + itb.Text[itb.KeyPos+1:]
	case ebiten.KeyEnd:
		itb.KeyPos = len(itb.Text)
	case ebiten.KeyHome:
		itb.KeyPos = 0
	case ebiten.KeyArrowLeft:
		if itb.KeyPos == 0 {
			break
		}
		itb.KeyPos--
	case ebiten.KeyArrowRight:
		if itb.KeyPos == len(itb.Text) {
			break
		}
		itb.KeyPos++
	case ebiten.KeyArrowUp:
	case ebiten.KeyArrowDown:
	default:
		itb.Text = itb.Text[:itb.KeyPos] + keyText + itb.Text[itb.KeyPos:]
		if key == ebiten.KeyTab {
			itb.KeyPos += 2
		} else {
			itb.KeyPos++
		}
	}
}
