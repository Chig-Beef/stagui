package stagui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	X, Y, W, H float64

	// Used to identify the button
	Name string

	// Whether the button can be used
	Disabled bool

	// Good for toggle buttons
	Pressed bool

	// Text to display on the button
	Text      string
	FontSize  float64
	TextColor color.Color

	// Color of image of button
	BgColor color.Color
	BgImg   *ebiten.Image

	// The color if pressed
	PressedColor color.Color

	// The color if disabled
	DisabledColor color.Color
}

func (b *Button) Draw(screen *ebiten.Image, ih ImageHandler, fh FontHandler) {
	// Draw bg
	if b.BgImg == nil {
		b.drawAsSolidColor(screen)
	} else {
		ih.DrawImage(screen, b.BgImg, float64(b.X), float64(b.Y), float64(b.W), float64(b.H), &ebiten.DrawImageOptions{})
	}

	// Draw text
	if b.Text != "" {
		op := text.DrawOptions{}
		op.PrimaryAlign = text.AlignCenter
		op.ColorScale.ScaleWithColor(b.TextColor)
		fh.DrawText(screen, b.Text, b.FontSize, float64(b.X+b.W/2), float64(b.Y), fh.GetFont("button"), &op)
	}
}

func (b *Button) drawAsSolidColor(screen *ebiten.Image) {
	if b.Disabled {
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), b.DisabledColor, false)
		return
	}

	if b.Pressed {
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), b.PressedColor, false)
	} else {
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), b.BgColor, false)
	}
}

func (b Button) CheckClick(x, y float64) bool {
	return !b.Disabled &&
		b.X <= x && x <= b.X+b.W &&
		b.Y <= y && y <= b.Y+b.H
}
