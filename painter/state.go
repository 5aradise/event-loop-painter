package painter

import (
	"image"
	"image/color"
	"slices"

	"github.com/roman-mazur/architecture-lab-3/ui"
	"golang.org/x/exp/shiny/screen"
)

type textureState struct {
	background struct {
		color color.Color
		rect  *Rectangle
	}
	figures []Point
}

func newTextureState() *textureState {
	s := &textureState{}
	s.background.color = color.Black
	return s
}

func (s *textureState) set(t screen.Texture) {
	t.Fill(t.Bounds(), s.background.color, screen.Src)
	if s.background.rect != nil {
		t.Fill(s.background.rect.
			Resize(image.Pt(t.Bounds().Dx(), t.Bounds().Dy())).
			ToImage().
			Add(t.Bounds().Min), color.Black, screen.Src)
	}
	for _, figure := range s.figures {
		ui.Figure(t, figure.
			Resize(image.Pt(t.Bounds().Dx(), t.Bounds().Dy())).
			ToImage().
			Add(t.Bounds().Min))
	}
}

func (s1 textureState) Equal(s2 textureState) bool {
	var isRectsEqual bool
	rect1, rect2 := s1.background.rect, s2.background.rect
	if rect1 != nil {
		if rect2 != nil {
			isRectsEqual = *rect1 == *rect2
		}
	} else {
		isRectsEqual = rect1 == rect2
	}

	return isColorsEqual(s1.background.color, s2.background.color) &&
		isRectsEqual &&
		slices.Equal(s1.figures, s2.figures)
}

func isColorsEqual(c1, c2 color.Color) bool {
	if c1 != nil {
		if c2 != nil {
			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, a2 := c2.RGBA()
			return r1 == r2 &&
				g1 == g2 &&
				b1 == b2 &&
				a1 == a2
		}
	} else {
		if c2 == nil {
			return true
		}
	}
	return false
}

func MockState() textureState {
	return textureState{figures: []Point{{0, 0}}}
}
