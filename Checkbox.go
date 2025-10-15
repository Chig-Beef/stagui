package stagui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Checkbox struct {
	// Used to identify the button
	Name string

	RadioGroups []string

	X, Y, W, H float64

	Bg, DisabledBg, CheckedBg ImageData

	Checked bool

	// Whether the button can be used
	Disabled bool

	// Good for toggle buttons
	Pressed bool
}

func (c *Checkbox) Draw(screen *ebiten.Image, vh VisualHandler) {
	if c.Disabled {
		c.DisabledBg.Draw(screen, vh, c.X, c.Y, c.W, c.H)
		return
	}

	if c.Checked {
		c.CheckedBg.Draw(screen, vh, c.X, c.Y, c.W, c.H)
		return
	}

	c.Bg.Draw(screen, vh, c.X, c.Y, c.W, c.H)
}

func (c *Checkbox) CheckFlip(x, y float64) bool {
	// TODO: Radio logic

	if c.Disabled {
		return false
	}

	held := c.X <= x && x <= c.X+c.W &&
		c.Y <= y && y <= c.Y+c.H

	if held && !c.Pressed {
		c.Pressed = true
		return true
	} else if !held && c.Pressed {
		c.Pressed = false
		return true
	}

	return false
}
