package stagui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Slider struct {
	Name    string // Used to identify the slider
	X, Y, W, H float64
	KnobW, KnobH float64
	Value float64

	LineImage, PressedLineImage, DisabledLineImage ImageData
	KnobImage, PressedKnobImage, DisabledKnobImage ImageData

	Pressed bool
	Disabled bool
	Vertical bool

	// Does this scroll the entire page?
	Pager bool
}

func (s *Slider) Draw(screen *ebiten.Image, vh VisualHandler) {
	s.drawLine(screen, vh)
	s.drawKnob(screen, vh)
}

func (s *Slider) drawLine(screen *ebiten.Image, vh VisualHandler) {
	if s.Disabled {
		s.DisabledLineImage.Draw(screen, vh, s.X, s.Y, s.W, s.H)
		return
	}

	if s.Pressed {
		s.PressedLineImage.Draw(screen, vh, s.X, s.Y, s.W, s.H)
		return
	}

	s.LineImage.Draw(screen, vh, s.X, s.Y, s.W, s.H)
}

func (s *Slider) getKnobPosDim() (float64, float64, float64, float64) {
	// TODO: Implement get knob pos dim
	return 0, 0, 0, 0
}

func (s *Slider) drawKnob(screen *ebiten.Image, vh VisualHandler) {
	x, y, w, h := s.getKnobPosDim()

	// TODO: Handle vertical cases

	if s.Disabled {
		s.DisabledKnobImage.Draw(screen, vh, x, y, w, h)
		return
	}

	if s.Pressed {
		s.PressedKnobImage.Draw(screen, vh, x, y, w, h)
		return
	}

	s.KnobImage.Draw(screen, vh, x, y, w, h)
}

func (s *Slider) CheckCollide(x, y float64) bool {
	return s.X <= x && x <= s.X+s.W &&
		s.Y <= y && y <= s.Y+s.H
}

func (s *Slider) Update(curMousePos [2]float64) bool {
	x := curMousePos[0]
	y := curMousePos[1]

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		return false
	}

	if !s.CheckCollide(x, y) {
		return false
	}

	dx := x - s.X

	s.Value = float64(dx / (s.W - 8))
	if s.Value > 1 {
		s.Value = 1
	}

	// Changed
	return true
}
