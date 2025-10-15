package stagui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type VisualHandler interface {
	DrawImage(*ebiten.Image, float64, float64, float64, float64, *ebiten.DrawImageOptions)
	DrawRect(float64, float64, float64, float64, color.Color)
	DrawText(string, float64, float64, float64, *text.GoTextFaceSource, *text.DrawOptions)
	GetFont(string) *text.GoTextFaceSource
	DrawLine(float64, float64, float64, float64, float64, color.Color)
	Translate(float64) float64
	GetImage(string) *ebiten.Image
}
