package main

import "fmt"

func main() {
	asiaBuilder := getBuilder("asia")
	europeBuilder := getBuilder("europe")

	god1 := NewDirector(europeBuilder)
	european := god1.CreatePerson()

	god2 := NewDirector(asiaBuilder)
	asian := god2.CreatePerson()
	fmt.Printf("European: %v\nAsian: %v\n", european, asian)
}
