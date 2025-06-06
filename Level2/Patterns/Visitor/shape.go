package main

import (
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
	accept(visitor Visitor)
}

type Rectangle struct {
	length, width float64
}

func (r Rectangle) Area() float64 {
	return r.length * r.width
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.length + r.width)
}

func (r Rectangle) accept(visitor Visitor) {
	visitor.ExtractRectangle(r)
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func (r Circle) accept(visitor Visitor) {
	visitor.ExtractCircle(r)
}

type Triangle struct {
	height, lowSide, leftSide, rightSide float64
}

func (t Triangle) Area() float64 {
	return t.height * t.lowSide / 2
}

func (t Triangle) Perimeter() float64 {
	return t.leftSide + t.rightSide + t.lowSide
}

func (r Triangle) accept(visitor Visitor) {
	visitor.ExtractTriangle(r)
}
