package painter

import (
	"image/color"
)

type Operation interface {
	Do(state *textureState) (ready bool)
}

type OperationList []Operation

func (ol OperationList) Do(t *textureState) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

type updateOp struct{}

func (op updateOp) Do(s *textureState) bool { return true }

var Update = updateOp{}

type OperationFunc func(s *textureState)

func (f OperationFunc) Do(s *textureState) bool {
	f(s)
	return false
}

var WhiteFill OperationFunc = func(s *textureState) {
	s.background.color = color.White
}

var GreenFill OperationFunc = func(s *textureState) {
	s.background.color = color.RGBA{G: 0xff, A: 0xff}
}

func BgRect(coords Rectangle) OperationFunc {
	return func(s *textureState) {
		s.background.rects = append(s.background.rects, coords)
	}
}

func Figure(coords Point) OperationFunc {
	return func(s *textureState) {
		s.figures = append(s.figures, coords)
	}
}

func Move(coords Point) OperationFunc {
	return func(s *textureState) {
		for i := range s.figures {
			s.figures[i] = coords
		}
	}
}

var Reset OperationFunc = func(s *textureState) {
	s.background.color = color.Black
	s.background.rects = s.background.rects[:0]
	s.figures = s.figures[:0]
}
