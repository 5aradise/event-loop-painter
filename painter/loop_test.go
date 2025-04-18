package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

type mockReceiver struct {
	textures []screen.Texture
}

func (tr *mockReceiver) Update(t screen.Texture) {
	tr.textures = append(tr.textures, t)
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	bgColor color.Color
	rects   []image.Rectangle
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	if dr == m.Bounds() {
		m.bgColor = src
	} else {
		m.rects = append(m.rects, dr)
	}
}

func TestBackgroundColoring(t *testing.T) {
	var (
		l = NewLoop()
		r mockReceiver
	)
	l.Receiver = &r
	go l.Start(mockScreen{})

	l.Post(GreenFill)
	l.Post(GreenFill)
	l.Post(GreenFill)
	l.Post(WhiteFill)
	l.Post(Update)
	l.StopAndWait()

	texture := r.textures[0].(*mockTexture)
	if !isColorsEqual(texture.bgColor, color.White) {
		t.Errorf("background is not white")
	}
}

func TestLastBgRect(t *testing.T) {
	var (
		l = NewLoop()
		r mockReceiver

		rect1 = Rect(1.1, 1.2, 1.3, 1.4)
		rect2 = Rect(2.1, 2.2, 2.3, 2.4)
		rect3 = Rect(3.1, 3.2, 3.3, 3.4)
		rect4 = Rect(4.1, 4.2, 4.3, 4.4)
	)
	l.Receiver = &r
	go l.Start(mockScreen{})

	l.Post(BgRect(rect1))
	l.Post(BgRect(rect2))
	l.Post(BgRect(rect3))
	l.Post(BgRect(rect4))

	l.Post(Update)
	l.StopAndWait()

	want := rect4.Resize(size).ToImage()

	texture := r.textures[0].(*mockTexture)
	if len(texture.rects) != 1 {
		t.Errorf("saved more than last BgRect")
	}
	if texture.rects[0] != want {
		t.Errorf("have: %v, want: %v", texture.rects[0], want)
	}
}

func TestFigure(t *testing.T) {
	var (
		l = NewLoop()
		r mockReceiver
	)
	l.Receiver = &r
	go l.Start(mockScreen{})

	l.Post(Figure(Point{}))

	l.Post(Update)
	l.StopAndWait()

	texture := r.textures[0].(*mockTexture)
	if len(texture.rects) == 0 {
		t.Errorf("figure does not create any shapes")
	}
}

func TestRectAndFigure(t *testing.T) {
	var (
		l = NewLoop()
		r mockReceiver

		rect = Rect(1, 2, 3, 4)
	)
	l.Receiver = &r
	go l.Start(mockScreen{})

	l.Post(BgRect(rect))
	l.Post(Figure(Point{0.5, 0.5}))
	l.Post(GreenFill)

	l.Post(Update)
	l.StopAndWait()

	texture := r.textures[0].(*mockTexture)
	if len(texture.rects) < 2 {
		t.Errorf("no square or shape was created")
	}
	if texture.rects[0] != rect.Resize(size).ToImage() {
		t.Errorf("wrong rect: have: %v, want: %v", texture.rects[0], rect.Resize(size).ToImage())
	}
	if !isColorsEqual(texture.bgColor, color.RGBA{G: 0xff, A: 0xff}) {
		t.Errorf("background color did not change to green")
	}
}

func TestReset(t *testing.T) {
	var (
		l = NewLoop()
		r mockReceiver
	)
	l.Receiver = &r
	go l.Start(mockScreen{})

	l.Post(Figure(Point{}))
	l.Post(BgRect(Rectangle{}))
	l.Post(WhiteFill)
	l.Post(Reset)

	l.Post(Update)
	l.StopAndWait()

	texture := r.textures[0].(*mockTexture)
	if len(texture.rects) != 0 {
		t.Errorf("texture contain shapes")
	}
	if !isColorsEqual(texture.bgColor, color.Black) {
		t.Errorf("texture is not black")
	}
}

func TestManyUpdates(t *testing.T) {
	var (
		l = NewLoop()
		r mockReceiver

		rect1 = Rect(1.1, 1.2, 1.3, 1.4)
		rect2 = Rect(2.1, 2.2, 2.3, 2.4)
		rect3 = Rect(3.1, 3.2, 3.3, 3.4)
	)
	l.Receiver = &r
	go l.Start(mockScreen{})

	l.Post(BgRect(rect1))
	l.Post(Update)

	l.Post(BgRect(rect2))
	l.Post(Update)

	l.Post(BgRect(rect3))
	l.Post(Update)
	l.StopAndWait()

	if len(r.textures) != 3 {
		t.Errorf("must be 3 separate textures, have: %d", len(r.textures))
	}

	texture := r.textures[0].(*mockTexture)
	if len(texture.rects) != 1 {
		t.Errorf("first texture has more than one shape, have: %d", len(texture.rects))
	}
	if texture.rects[0] != rect1.Resize(size).ToImage() {
		t.Errorf("first texture does not have wanted shape, have: %v, want: %v", texture.rects[0], rect1.Resize(size).ToImage())
	}

	texture = r.textures[1].(*mockTexture)
	if len(texture.rects) != 1 {
		t.Errorf("second texture has more than one shape, have: %d", len(texture.rects))
	}
	if texture.rects[0] != rect2.Resize(size).ToImage() {
		t.Errorf("second texture does not have wanted shape, have: %v, want: %v", texture.rects[0], rect2.Resize(size).ToImage())
	}

	texture = r.textures[2].(*mockTexture)
	if len(texture.rects) != 1 {
		t.Errorf("third texture has more than one shape, have: %d", len(texture.rects))
	}
	if texture.rects[0] != rect3.Resize(size).ToImage() {
		t.Errorf("third texture does not have wanted shape, have: %v, want: %v", texture.rects[0], rect3.Resize(size).ToImage())
	}
}

func TestStopAndWaint(t *testing.T) {
	var (
		l = NewLoop()
		r mockReceiver
	)
	l.Receiver = &r
	go l.Start(mockScreen{})

	l.StopAndWait()

	l.Post(GreenFill)
	l.Post(Update)

	if len(r.textures) != 0 {
		t.Errorf("update has an effect after closing")
	}
}

func TestMove(t *testing.T) {
	var (
		l = NewLoop()
		r mockReceiver

		rect = Rect(1, 1, 1, 1)
		move = Pt(2, 2)
	)
	l.Receiver = &r
	go l.Start(mockScreen{})

	l.Post(BgRect(rect))
	l.Post(Move(move))
	l.Post(Update)

	l.StopAndWait()

	texture := r.textures[0].(*mockTexture)
	if texture.rects[0] != rect.Resize(size).ToImage() {
		t.Errorf("move has an effect on BgRect")
	}
}
