package main

import "fmt"

func main() {
	a := 0b1111

	fmt.Printf("a = %b\n", a&^0b0100)
}
