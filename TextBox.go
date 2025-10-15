package stagui

import (
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type TextBox struct {
	Name string
	X, Y, W, H float64
	Text      TextData
	Bg, ActiveBg, DisabledBg ImageData

	// Is the user attempting to type
	// inside the textbox
	Active bool

	Disabled bool

	// Limit to a single line
	Inline bool

	// Where in the textbox the user is
	// typing
	KeyPosX, KeyPosY int

	// Prevents bleeding
	buffer *ebiten.Image
}

func (tb *TextBox) Draw(screen *ebiten.Image, vh VisualHandler) {
	if tb.checkRecreateBuffer() {
		// TODO: I don't think this is correct size
		tb.buffer = ebiten.NewImage(int(tb.W), int(tb.H))
	}

	// TODO: Draw each individual text
	tb.drawBackground(screen, vh)
	tb.Text.draw(tb.buffer, vh, 0, 0)
	vh.DrawImage(tb.buffer, tb.X, tb.Y, tb.W, tb.H, &ebiten.DrawImageOptions{})

	// Draw a line at the bottom of the textbox
	if tb.Active {
		vh.DrawLine(tb.X, tb.Y+tb.H, tb.X+tb.W, tb.Y, 1, color.White)
	}
}

func (tb *TextBox) checkRecreateBuffer() bool {
	if tb.buffer == nil {
		return true
	}

	// TODO: Correct dimensions

	return false
}

func (tb *TextBox) drawBackground(screen *ebiten.Image, vh VisualHandler) {
	if tb.Disabled {
		tb.DisabledBg.Draw(screen, vh, tb.X, tb.Y, tb.W, tb.H)
		return
	}

	if tb.Active {
		tb.ActiveBg.Draw(screen, vh, tb.X, tb.Y, tb.W, tb.H)
		return
	}

	tb.Bg.Draw(screen, vh, tb.X, tb.Y, tb.W, tb.H)
}

func (tb *TextBox) CheckClick(x, y float64) bool {
	clicked := tb.X <= x && x <= tb.X+tb.W &&
		tb.Y <= y && y <= tb.Y+tb.H

	// Effectively an xor operation.
	// If we're active and clicked, set
	// active to false
	tb.Active = clicked != tb.Active

	return clicked
}

func (tb *TextBox) addCharacter(char string) {
	tb.Text.Text[tb.KeyPosY] = tb.Text.Text[tb.KeyPosY][:tb.KeyPosX] + char + tb.Text.Text[tb.KeyPosY][tb.KeyPosX:]

	// To accomodate longer strings such as tab "  "
	tb.KeyPosX += len(char)
}

func (tb *TextBox) move(dir ebiten.Key) {
	switch dir {
	case ebiten.KeyArrowLeft:
		tb.KeyPosX--

		if tb.KeyPosX >= 0 {
			break
		}

		if tb.Inline {
			tb.KeyPosX = 0
		} else {
			tb.KeyPosY--
			if tb.KeyPosY < 0 {
				tb.KeyPosY = 0
				tb.KeyPosX = 0
			} else {
				if tb.KeyPosX > len(tb.Text.Text[tb.KeyPosY]) {
					tb.KeyPosX = len(tb.Text.Text[tb.KeyPosY])
				}
			}
		}

	case ebiten.KeyArrowRight:
		tb.KeyPosX++

		if tb.KeyPosX <= len(tb.Text.Text[tb.KeyPosY]) {
			break
		}

		if tb.Inline {
			tb.KeyPosX = len(tb.Text.Text[tb.KeyPosY])
		} else {
			tb.KeyPosY++
			if tb.KeyPosY >= len(tb.Text.Text) {
				tb.KeyPosY = len(tb.Text.Text)-1
				tb.KeyPosX = len(tb.Text.Text[tb.KeyPosY])
			} else {
				tb.KeyPosX = 0
			}
		}

	case ebiten.KeyArrowUp:
		if tb.Inline {
			break
		}

		tb.KeyPosY--
		if tb.KeyPosY < 0 {
			tb.KeyPosY = 0
		}

		if tb.KeyPosX > len(tb.Text.Text[tb.KeyPosY]) {
			tb.KeyPosX = len(tb.Text.Text[tb.KeyPosY])
		}

	case ebiten.KeyArrowDown:
		if tb.Inline {
			break
		}

		tb.KeyPosY++
		if tb.KeyPosY >= len(tb.Text.Text) {
			tb.KeyPosY = len(tb.Text.Text) - 1
		}

		if tb.KeyPosX > len(tb.Text.Text[tb.KeyPosY]) {
			tb.KeyPosX = len(tb.Text.Text[tb.KeyPosY])
		}

	case ebiten.KeyEnd:
		tb.KeyPosX = len(tb.Text.Text[tb.KeyPosY])

	case ebiten.KeyHome:
		tb.KeyPosX = 0
	}
}

func (tb *TextBox) Update() {
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
		if tb.KeyPosX == len(tb.Text.Text[tb.KeyPosY]) {
			tb.Text.Text = slices.Insert(tb.Text.Text, tb.KeyPosY+1, "")
			tb.KeyPosY++
			tb.KeyPosX = 0
			break
		}

		tb.Text.Text = slices.Insert(tb.Text.Text, tb.KeyPosY+1, tb.Text.Text[tb.KeyPosY][tb.KeyPosX:])
		tb.Text.Text[tb.KeyPosY] = tb.Text.Text[tb.KeyPosY][:tb.KeyPosX]

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

			tb.Text.Text[tb.KeyPosY] = tb.Text.Text[tb.KeyPosY][:tb.KeyPosX-1] + tb.Text.Text[tb.KeyPosY][tb.KeyPosX:]
			tb.KeyPosX--
			break
		}

		// Any other line

		// Right at the start
		if tb.KeyPosX == 0 {
			// Add this line to the previous line
			tb.Text.Text[tb.KeyPosY-1] += tb.Text.Text[tb.KeyPosY]

			// Delete the line
			tb.Text.Text = slices.Delete(tb.Text.Text, tb.KeyPosY, tb.KeyPosY+1)

			tb.KeyPosY--
			tb.KeyPosX = len(tb.Text.Text[tb.KeyPosY])
			break
		}

		// Take bit out of line
		tb.Text.Text[tb.KeyPosY] = tb.Text.Text[tb.KeyPosY][:tb.KeyPosX-1] + tb.Text.Text[tb.KeyPosY][tb.KeyPosX:]
		tb.KeyPosX--

	case ebiten.KeyDelete:
		// Complete end of text
		if tb.KeyPosY == len(tb.Text.Text)-1 {
			// End of line
			if tb.KeyPosX == len(tb.Text.Text[tb.KeyPosY]) {
				// Does nothing
				break
			}

			// Delete the character
			tb.Text.Text[tb.KeyPosY] = tb.Text.Text[tb.KeyPosY][:tb.KeyPosX] + tb.Text.Text[tb.KeyPosY][tb.KeyPosX+1:]
			break
		}

		// Somewhere beforehand
		if tb.KeyPosX == len(tb.Text.Text[tb.KeyPosY]) {
			// Move next line onto this line
			tb.Text.Text[tb.KeyPosY] += tb.Text.Text[tb.KeyPosY+1]

			// Delete the next line
			tb.Text.Text = slices.Delete(tb.Text.Text, tb.KeyPosY+1, tb.KeyPosY+2)
			break
		}

		// There's no character to delete
		if len(tb.Text.Text[tb.KeyPosY]) == 0 {
			break
		}

		// Delete that single character
		tb.Text.Text[tb.KeyPosY] = tb.Text.Text[tb.KeyPosY][:tb.KeyPosX] + tb.Text.Text[tb.KeyPosY][tb.KeyPosX+1:]

	case ebiten.KeyArrowLeft,
		ebiten.KeyArrowRight,
		ebiten.KeyArrowUp,
		ebiten.KeyArrowDown,
		ebiten.KeyEnd,
		ebiten.KeyHome:
		tb.move(key)

	default:
		tb.addCharacter(keyText)
	}
}
