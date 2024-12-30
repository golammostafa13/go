package main

import (
    "fmt"
    "go_tutorials/cmd/tutorial_1" // Import the tutorial_1 module from the correct path
)

func main() {
	fmt.Println("Hello from go_tutorial_2!")
	tutorial_1.PrintMessage() // Call a function from tutorial_1 module
}
