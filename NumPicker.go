package stagui

import (
	"strconv"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type NumPicker struct {
	Name string
	X, Y, W, H float64
	Text      TextData
	Bg, ActiveBg, DisabledBg ImageData

	AllowFloat bool
	IntData int
	FloatData float64

	Min, Max float64

	// Is the user attempting to type
	// inside the textbox
	Active bool

	Disabled bool

	// Where in the textbox the user is
	// typing
	KeyPosX int

	// Prevents bleeding
	buffer *ebiten.Image
}

func (np *NumPicker) Draw(screen *ebiten.Image, vh VisualHandler) {
	if np.checkRecreateBuffer() {
		// TODO: I don't think this is correct size
		np.buffer = ebiten.NewImage(int(np.W), int(np.H))
	}

	// TODO: Draw each individual text
	np.drawBackground(screen, vh)
	np.Text.draw(np.buffer, vh, 0, 0)
	vh.DrawImage(np.buffer, np.X, np.Y, np.W, np.H, &ebiten.DrawImageOptions{})

	// Draw a line at the bottom of the num picker
	if np.Active {
		vh.DrawLine(np.X, np.Y+np.H, np.X+np.W, np.Y, 1, color.White)
	}
}

func (np *NumPicker) checkRecreateBuffer() bool {
	if np.buffer == nil {
		return true
	}

	// TODO: Correct dimensions

	return false
}

func (np *NumPicker) drawBackground(screen *ebiten.Image, vh VisualHandler) {
	if np.Disabled {
		np.DisabledBg.Draw(screen, vh, np.X, np.Y, np.W, np.H)
		return
	}

	if np.Active {
		np.ActiveBg.Draw(screen, vh, np.X, np.Y, np.W, np.H)
		return
	}

	np.Bg.Draw(screen, vh, np.X, np.Y, np.W, np.H)
}

func (np *NumPicker) CheckClick(x, y float64) bool {
	clicked := np.X <= x && x <= np.X+np.W &&
		np.Y <= y && y <= np.Y+np.H

	// Effectively an xor operation.
	// If we're active and clicked, set
	// active to false
	np.Active = clicked != np.Active

	np.saveValue()

	return clicked
}

func (np *NumPicker) saveValue() {
	// Ignore errors
	np.IntData, _ = strconv.Atoi(np.Text.Text[0])
	np.FloatData, _ = strconv.ParseFloat(np.Text.Text[0], 64)
}

func (np *NumPicker) addCharacter(char string) {
	np.Text.Text[0] = np.Text.Text[0][:np.KeyPosX] + char + np.Text.Text[0][np.KeyPosX:]

	// To accomodate longer strings such as tab "  "
	np.KeyPosX += len(char)
}

func (np *NumPicker) move(dir ebiten.Key) {
	switch dir {
	case ebiten.KeyArrowLeft:
		np.KeyPosX--

		if np.KeyPosX >= 0 {
			break
		}

		np.KeyPosX = 0

	case ebiten.KeyArrowRight:
		np.KeyPosX++

		if np.KeyPosX <= len(np.Text.Text[0]) {
			break
		}

		np.KeyPosX = len(np.Text.Text[0])

	case ebiten.KeyEnd:
		np.KeyPosX = len(np.Text.Text[0])

	case ebiten.KeyHome:
		np.KeyPosX = 0
	}
}

func (np *NumPicker) Update() {
	if !np.Active {
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
		np.Active = false
		np.saveValue()

	case ebiten.KeyBackspace:
		// Right at the start
		if np.KeyPosX == 0 {
			// Does nothing
			break
		}

		np.Text.Text[0] = np.Text.Text[0][:np.KeyPosX-1] + np.Text.Text[0][np.KeyPosX:]
		np.KeyPosX--

	case ebiten.KeyDelete:
			// End of line
			if np.KeyPosX == len(np.Text.Text[0]) {
				// Does nothing
				break
			}

			// Delete the character
			np.Text.Text[0] = np.Text.Text[0][:np.KeyPosX] + np.Text.Text[0][np.KeyPosX+1:]

	case ebiten.KeyArrowLeft,
		ebiten.KeyArrowRight,
		ebiten.KeyArrowUp,
		ebiten.KeyArrowDown,
		ebiten.KeyEnd,
		ebiten.KeyHome:
		np.move(key)

	default:
		np.addCharacter(keyText)
	}
}
