package main

import (
	"fmt"
	"math"
)

// Shape interface
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle type
type Rectangle struct {
	width, height float64
}

// Implementing Shape interface for Rectangle
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

// Circle type
type Circle struct {
	radius float64
}

// Implementing Shape interface for Circle
func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func main_interface() {
	// Create a slice of Shape interface
	shapes := []Shape{
		Rectangle{width: 5, height: 3},
		Circle{radius: 4},
	}

	// Iterate over the shapes and call the Area and Perimeter methods
	for _, shape := range shapes {
		fmt.Printf("Shape: %T\n", shape)
		fmt.Println("Area:", shape.Area())
		fmt.Println("Perimeter:", shape.Perimeter())
	}
}
