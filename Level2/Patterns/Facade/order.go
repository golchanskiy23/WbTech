package main

type Order struct {
	Items []Item
}

type Item struct {
	Name  string
	Value int
}
