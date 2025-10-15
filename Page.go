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
	Name string

	// Drawn at the top middle of the page
	Title string

	// Content of the page
	Buttons         []*Button
	Sliders         []*Slider
	Checkboxes      []*Checkbox
	Images          []*Image
	NumPickers      []*NumPicker
	TextBoxes       []*TextBox
	PlainTexts      []*PlainText
	Panels          []*Panel

	// Whether the bg will be drawn
	BgDraw bool

	// Color of the bg
	BgColor color.Color

	// Image for the bg
	BgImg *ebiten.Image
}

func (p *Page) GetButton(name string) *Button {
	for i := range len(p.Buttons) {
		if p.Buttons[i].Name == name {
			return p.Buttons[i]
		}
	}
	return nil
}

func (p *Page) GetSlider(name string) *Slider {
	for i := range len(p.Sliders) {
		if p.Sliders[i].Name == name {
			return p.Sliders[i]
		}
	}
	return nil
}

func (p *Page) GetText(name string) *PlainText {
	for i := range len(p.PlainTexts) {
		if p.PlainTexts[i].Name == name {
			return p.PlainTexts[i]
		}
	}
	return nil
}

func (p *Page) GetTextbox(name string) *TextBox {
	for i := range len(p.TextBoxes) {
		if p.TextBoxes[i].Name == name {
			return p.TextBoxes[i]
		}
	}
	return nil
}

func (p *Page) GetPanel(name string) *Panel {
	for i := range len(p.Panels) {
		if p.Panels[i].Name == name {
			return p.Panels[i]
		}
	}
	return nil
}

// Interaction logic of all content
func (p *Page) Update(curMousePos [2]float64) (string, *Button, *Slider) {
	for _, s := range p.Sliders {
		if s.Update(curMousePos) {
			return s.Name, nil, s
		}
	}

	// Key press
	for _, tb := range p.TextBoxes {
		tb.Update()
	}

	// Check whether they're even pressing
	// the left mouse button. We don't care
	// about any other button press
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		return "", nil, nil
	}

	for _, b := range p.Buttons {
		if b.CheckClick(curMousePos[0], curMousePos[1]) {
			return b.Name, b, nil
		}
	}

	// Mouse press
	for _, tb := range p.TextBoxes {
		tb.CheckClick(curMousePos[0], curMousePos[1])
	}

	return "", nil, nil
}

func (p *Page) Draw(screen *ebiten.Image, vh VisualHandler) {
	// Really I don't like that we're
	// filling the screen every frame for
	// no apparent reason, but it's a title
	// screen, so I don't know why I care
	if p.BgDraw {
		if p.BgImg != nil {
			vh.DrawImage(p.BgImg, 0, 0, 100, 50, &ebiten.DrawImageOptions{})
		} else {
			screen.Fill(p.BgColor)
		}
	}

	// Draw the title
	if p.Title != "" {
		op := text.DrawOptions{}
		op.PrimaryAlign = text.AlignCenter
		vh.DrawText(p.Title, 8, 50, 1, vh.GetFont("default"), &op)
	}

	// Content

	for _, t := range p.PlainTexts {
		t.Draw(screen, vh)
	}

	for _, b := range p.Buttons {
		b.Draw(screen, vh)
	}

	for _, s := range p.Sliders {
		s.Draw(screen, vh)
	}

	for _, tb := range p.TextBoxes {
		tb.Draw(screen, vh)
	}

	for _, panel := range p.Panels {
		panel.Draw(screen, vh)
	}
}
