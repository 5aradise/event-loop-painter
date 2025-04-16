package painter

import "image"

type Point struct {
	X, Y float32
}

func Pt(x, y float32) Point {
	return Point{x, y}
}

func (p Point) Resize(d image.Point) Point {
	return Point{p.X * float32(d.X), p.Y * float32(d.Y)}
}

func (p Point) ToImage() image.Point {
	return image.Pt(int(p.X), int(p.Y))
}

type Rectangle struct {
	Min, Max Point
}

func Rect(x0, y0, x1, y1 float32) Rectangle {
	return Rectangle{Point{x0, y0}, Point{x1, y1}}
}

func (r Rectangle) Resize(d image.Point) Rectangle {
	r.Min = r.Min.Resize(d)
	r.Max = r.Max.Resize(d)
	return r
}

func (r Rectangle) ToImage() image.Rectangle {
	return image.Rectangle{
		Min: r.Min.ToImage(),
		Max: r.Max.ToImage(),
	}
}
