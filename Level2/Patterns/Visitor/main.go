package main

import "math"

func main() {
	arr := []Shape{Triangle{4, 4, math.Sqrt(20), math.Sqrt(20)}, Rectangle{5, 4}, Circle{3}}
	visitor := NewXMLExtractor()

	for _, shape := range arr {
		shape.accept(visitor)
	}
}
