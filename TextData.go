package stagui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Generic object for holding text for the screen
// e.g., used for buttons, text, text boxes, etc
type TextData struct {
	Text []string
	Color color.Color
	Font  *text.GoTextFaceSource
	Size  float64
	TextGap float64
	HorAlign, VerAlign text.Align
}

func (td *TextData) draw(screen *ebiten.Image, vh VisualHandler, x, y float64) {
	for i := range len(td.Text) {
		op := text.DrawOptions{}
		op.PrimaryAlign = td.HorAlign
		op.SecondaryAlign = td.VerAlign
		vh.DrawText(td.Text[i], td.Size, x, y+td.TextGap*float64(i), td.Font, &op)
	}
}
