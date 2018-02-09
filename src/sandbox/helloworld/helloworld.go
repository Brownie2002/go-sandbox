package main

import (
	"fmt"
	"sandbox/postgre"
)

func main() {
	fmt.Println("Hello world")

	fmt.Println(postgre.Reverse("Test"))
}
