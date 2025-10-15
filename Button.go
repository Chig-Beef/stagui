package stagui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Button struct {
	// Used to identify the button
	Name string

	X, Y, W, H float64

	Text, PressedText, DisabledText TextData
	Bg, PressedBg, DisabledBg ImageData

	// Whether the button can be used
	Disabled bool

	// Good for toggle buttons
	Pressed bool
}

func (b *Button) Draw(screen *ebiten.Image, vh VisualHandler) {
	b.drawBackground(screen, vh)
	b.Text.draw(screen, vh, b.getTextX(), b.getTextY())
}

func (b *Button) getTextX() float64 {
	switch b.Text.HorAlign {
	case text.AlignStart:
		return b.X
	case text.AlignCenter:
		return b.X+b.W/2
	case text.AlignEnd:
		return b.X+b.W
	}
	return b.X
}

// TODO: Check this works correctly
func (b *Button) getTextY() float64 {
	switch b.Text.VerAlign {
	case text.AlignStart:
		return b.Y
	case text.AlignCenter:
		return b.Y+b.H/2
	case text.AlignEnd:
		return b.Y+b.H
	}
	return b.Y
}

func (b *Button) drawBackground(screen *ebiten.Image, vh VisualHandler) {
	if b.Disabled {
		b.DisabledBg.Draw(screen, vh, b.X, b.Y, b.W, b.H)
		return
	}

	if b.Pressed {
		b.PressedBg.Draw(screen, vh, b.X, b.Y, b.W, b.H)
		return
	}

	b.Bg.Draw(screen, vh, b.X, b.Y, b.W, b.H)
}

func (b Button) CheckClick(x, y float64) bool {
	return !b.Disabled &&
		b.X <= x && x <= b.X+b.W &&
		b.Y <= y && y <= b.Y+b.H
}
