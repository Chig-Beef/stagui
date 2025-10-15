package stagui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Panel struct {
	Name string

	X, Y, W, H float64

	Bg ImageData
}

func (p *Panel) Draw(screen *ebiten.Image, vh VisualHandler) {
	p.Bg.Draw(screen, vh, p.X, p.Y, p.W, p.H)
}
