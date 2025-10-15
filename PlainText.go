package stagui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type PlainText struct {
	Name string
	Text TextData
	X float64
	Y float64
}

func (text *PlainText) Draw(screen *ebiten.Image, vh VisualHandler) {
	text.Text.draw(screen, vh, text.X, text.Y)
}
