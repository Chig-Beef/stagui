package stagui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Used to store an image or a color if none is available
type ImageData struct {
	Image *ebiten.Image
	Color color.Color
}

func (i *ImageData) Draw(screen *ebiten.Image, vh VisualHandler, x, y, w, h float64) {
	if i.Image == nil {
		vh.DrawRect(x, y, w, h, i.Color)
	} else {
		vh.DrawImage(i.Image, x, y, w, h, &ebiten.DrawImageOptions{})
	}
}
