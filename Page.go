package stagui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Page is for things like the menu,
// the settings page, etc.
type Page struct {
	// Drawn at the top middle of the page
	Title string

	// Content of the page
	Buttons         []*Button
	Sliders         []*Slider
	Text            []*StaticText
	InlineTextBoxes []*InlineTextBox

	// Whether the bg will be drawn
	BgDraw bool

	// Color of the bg
	BgColor color.Color

	// Image for the bg
	BgImg *ebiten.Image
}

// Interaction logic of all content
func (p *Page) Update(curMousePos [2]int) (string, *Button, *Slider) {
	for _, s := range p.Sliders {
		if s.Update(curMousePos) {
			return s.Name, nil, s
		}
	}

	// Key press
	for _, itb := range p.InlineTextBoxes {
		itb.Update()
	}

	// Check whether they're even pressing
	// the left mouse button. We don't care
	// about any other button press
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		return "", nil, nil
	}

	for _, b := range p.Buttons {
		if b.CheckClick(float64(curMousePos[0]), float64(curMousePos[1])) {
			return b.Name, b, nil
		}
	}

	// Mouse press
	for _, itb := range p.InlineTextBoxes {
		itb.CheckClick(curMousePos[0], curMousePos[1])
	}

	return "", nil, nil
}

func (p *Page) Draw(screen *ebiten.Image, ih ImageHandler, fh FontHandler, sw, sh, hw int) {
	// Really I don't like that we're
	// filling the screen every frame for
	// no apparent reason, but it's a title
	// screen, so I don't know why I care
	if p.BgDraw {
		if p.BgImg != nil {
			ih.DrawImage(screen, p.BgImg, 0, 0, float64(sw), float64(sh), &ebiten.DrawImageOptions{})
		} else {
			screen.Fill(p.BgColor)
		}
	}

	// Draw the title
	if p.Title != "" {
		op := text.DrawOptions{}
		op.PrimaryAlign = text.AlignCenter
		fh.DrawText(screen, p.Title, 100, float64(hw), 10, fh.GetFont("default"), &op)
	}

	// Content

	for _, t := range p.Text {
		t.Draw(screen)
	}

	for _, b := range p.Buttons {
		b.Draw(screen, ih, fh)
	}

	for _, s := range p.Sliders {
		s.Draw(screen)
	}

	for _, itb := range p.InlineTextBoxes {
		itb.Draw(screen, fh)
	}
}
