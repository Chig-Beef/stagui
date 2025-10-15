package stagui

import "github.com/hajimehoshi/ebiten/v2"

type Image struct {
	Name string
	X, Y, W, H float64
	Image ImageData
}

func (i *Image) Draw(screen *ebiten.Image, vh VisualHandler) {
	i.Image.Draw(screen, vh, i.X, i.Y, i.W, i.H)
}
