package utils

type Point struct {
  X float64
  Y float64
}

func NewPoint(x, y float64) Point {
  return Point{
    X: x,
    Y: y,
  }
}

func Pt(x, y float64) Point {
  return NewPoint(x, y)
}

