package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

const WindowSide = 800

var (
	Green  = color.RGBA{R: 100, G: 200, B: 100}
	Yellow = color.RGBA{R: 255, G: 200, B: 100}
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Point
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.X = WindowSide / 2
	pw.pos.Y = WindowSide / 2
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  WindowSide,
		Height: WindowSide,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Button == mouse.ButtonLeft {
				x := int(e.X)
				y := int(e.Y)
				pw.pos = image.Point{x, y}
				pw.w.Send(paint.Event{})
			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), Green, draw.Src) // Фон.

	pw.DrawT()

	// Малювання білої рамки.
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func drawTFigure(
	fillFunc func(r image.Rectangle, c color.Color, op draw.Op),
	pos image.Point,
	bounds image.Rectangle,
) {
	x := pos.X
	y := pos.Y

	width := bounds.Dx()
	height := bounds.Dy()

	horizW := width / 2
	horizH := height / 6
	vertW := width / 6
	vertH := height / 4

	rect1 := image.Rect(x-horizW/2, y, x+horizW/2, y+horizH)
	rect2 := image.Rect(x-vertW/2, y-vertH, x+vertW/2, y)

	fillFunc(rect1, Yellow, draw.Src)
	fillFunc(rect2, Yellow, draw.Src)
}

func (pw *Visualizer) DrawT() {
	drawTFigure(pw.w.Fill, pw.pos, pw.sz.Bounds())
}

func Figure(t screen.Texture, pos image.Point) {
	drawTFigure(t.Fill, pos, t.Bounds())
}
