package stagui

import (
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Like a text box, but only has one
// line
type TextBox struct {
	X, Y, W, H float64

	Text      []string
	TextColor color.Color
	FontSize  float64

	BgColor color.Color

	// Is the user attempting to type
	// inside the textbox
	Active bool

	// Where in the textbox the user is
	// typing
	KeyPosX, KeyPosY int
}

func (tb *TextBox) Draw(screen *ebiten.Image, fh FontHandler) {
	vector.DrawFilledRect(screen, float32(tb.X), float32(tb.Y), float32(tb.W), float32(tb.H), tb.BgColor, false)

	for i := range len(tb.Text) {
		fh.DrawText(screen, tb.Text[i], tb.FontSize, tb.X+4, tb.Y+2+float64(i)*(tb.FontSize+2), fh.GetFont("textBox"), &text.DrawOptions{})
	}

	// Draw a line at the bottom of the textbox
	if tb.Active {
		vector.StrokeLine(screen, float32(tb.X), float32(tb.Y+tb.H), float32(tb.X+tb.W), float32(tb.Y+tb.H), 1, color.White, false)
	}
}

func (tb *TextBox) CheckClick(mx, my int) bool {
	x := float64(mx)
	y := float64(my)

	clicked := tb.X <= x && x <= tb.X+tb.W &&
		tb.Y <= y && y <= tb.Y+tb.H

	// Effectively an xor operation.
	// If we're active and clicked, set
	// active to false
	tb.Active = clicked != tb.Active

	return clicked
}

func (tb *TextBox) Update() {
	// ASSERT len(tb.Text) != 0

	if !tb.Active {
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
		if tb.KeyPosX == len(tb.Text[tb.KeyPosY]) {
			tb.Text = slices.Insert(tb.Text, tb.KeyPosY+1, "")
			tb.KeyPosY++
			tb.KeyPosX = 0
			break
		}

		tb.Text = slices.Insert(tb.Text, tb.KeyPosY+1, tb.Text[tb.KeyPosY][tb.KeyPosX:])
		tb.Text[tb.KeyPosY] = tb.Text[tb.KeyPosY][:tb.KeyPosX]

		tb.KeyPosY++
		tb.KeyPosX = 0

	case ebiten.KeyBackspace:
		// First line
		if tb.KeyPosY == 0 {
			// Right at the start
			if tb.KeyPosX == 0 {
				// Doest nothing
				break
			}

			tb.Text[tb.KeyPosY] = tb.Text[tb.KeyPosY][:tb.KeyPosX-1] + tb.Text[tb.KeyPosY][tb.KeyPosX:]
			tb.KeyPosX--
			break
		}

		// Any other line

		// Right at the start
		if tb.KeyPosX == 0 {
			// Add this line to the previous line
			tb.Text[tb.KeyPosY-1] += tb.Text[tb.KeyPosY]

			// Delete the line
			tb.Text = slices.Delete(tb.Text, tb.KeyPosY, tb.KeyPosY+1)

			tb.KeyPosY--
			tb.KeyPosX = len(tb.Text[tb.KeyPosY])
			break
		}

		// Take bit out of line
		tb.Text[tb.KeyPosY] = tb.Text[tb.KeyPosY][:tb.KeyPosX-1] + tb.Text[tb.KeyPosY][tb.KeyPosX:]
		tb.KeyPosX--

	case ebiten.KeyDelete:
		// Complete end of text
		if tb.KeyPosY == len(tb.Text)-1 {
			// End of line
			if tb.KeyPosX == len(tb.Text[tb.KeyPosY]) {
				// Does nothing
				break
			}

			// Delete the character
			tb.Text[tb.KeyPosY] = tb.Text[tb.KeyPosY][:tb.KeyPosX] + tb.Text[tb.KeyPosY][tb.KeyPosX+1:]
			break
		}

		// Somewhere beforehand
		if tb.KeyPosX == len(tb.Text[tb.KeyPosY]) {
			// Move next line onto this line
			tb.Text[tb.KeyPosY] += tb.Text[tb.KeyPosY+1]

			// Delete the next line
			tb.Text = slices.Delete(tb.Text, tb.KeyPosY+1, tb.KeyPosY+2)
			break
		}

		// There's no character to delete
		if len(tb.Text[tb.KeyPosY]) == 0 {
			break
		}

		// Delete that single character
		tb.Text[tb.KeyPosY] = tb.Text[tb.KeyPosY][:tb.KeyPosX] + tb.Text[tb.KeyPosY][tb.KeyPosX+1:]

	case ebiten.KeyEnd:
		tb.KeyPosX = len(tb.Text[tb.KeyPosY])

	case ebiten.KeyHome:
		tb.KeyPosX = 0

	case ebiten.KeyArrowLeft:
		if tb.KeyPosX == 0 {
			tb.KeyPosY--
			if tb.KeyPosY < 0 {
				tb.KeyPosY = 0
			}
			break
		}
		tb.KeyPosX--

	case ebiten.KeyArrowRight:
		if tb.KeyPosX == len(tb.Text[tb.KeyPosX]) {
			if tb.KeyPosY < len(tb.Text) {
				tb.KeyPosY++
			}
			break
		}
		tb.KeyPosX++

	case ebiten.KeyArrowUp:
		tb.KeyPosY--
		if tb.KeyPosY < 0 {
			tb.KeyPosY = 0
		}

		if tb.KeyPosX > len(tb.Text[tb.KeyPosY]) {
			tb.KeyPosX = len(tb.Text[tb.KeyPosY])
		}

	case ebiten.KeyArrowDown:
		tb.KeyPosY++
		if tb.KeyPosY >= len(tb.Text) {
			tb.KeyPosY = len(tb.Text) - 1
		}

		if tb.KeyPosX > len(tb.Text[tb.KeyPosY]) {
			tb.KeyPosX = len(tb.Text[tb.KeyPosY])
		}

	default:
		tb.Text[tb.KeyPosY] = tb.Text[tb.KeyPosY][:tb.KeyPosX] + keyText + tb.Text[tb.KeyPosY][tb.KeyPosX:]
		if key == ebiten.KeyTab {
			tb.KeyPosX += 2
		} else {
			tb.KeyPosX++
		}
	}
}
