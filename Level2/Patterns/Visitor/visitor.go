package main

import "fmt"

type Visitor interface {
	ExtractTriangle(Triangle)
	ExtractRectangle(Rectangle)
	ExtractCircle(Circle)
}

type XMLExtractor struct {
}

func NewXMLExtractor() *XMLExtractor {
	return &XMLExtractor{}
}

func (e *XMLExtractor) ExtractTriangle(triangle Triangle) {
	fmt.Printf("Extrction triangle : %v to xml format completed!!!\n", triangle)
}

func (e *XMLExtractor) ExtractRectangle(rectangle Rectangle) {
	fmt.Printf("Extrction rectangle : %v to xml format completed!!!\n", rectangle)
}

func (e *XMLExtractor) ExtractCircle(circle Circle) {
	fmt.Printf("Extrction circle : %v to xml format completed!!!\n", circle)
}
