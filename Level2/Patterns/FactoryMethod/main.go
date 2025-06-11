package main

import "fmt"

type Transport struct {
	PassengerCapacity int
	MaximumLoad       int
	Fuel              float32
	EnginePower       float32
}

type Deliver interface {
	Delivery(source, destination string, load int)
}

type Ship struct {
	SeaMileage float32
	Transport
}

func (s *Transport) Delivery(source, destination string, load int) {
	fmt.Printf("...", source, destination, load)
}

func (s *Ship) Delivery(source, destination string, load int) {
	fmt.Printf("Ship delivering from %s to %s with load %d\n", source, destination, load)
}

type Truck struct {
	TiresAmount int
	Transport
}

func (t *Truck) Delivery(source, destination string, load int) {
	fmt.Printf("Truck delivering from %s to %s with load %d\n", source, destination, load)
}

type Order struct {
	Name string
	Load int
}

type DeliverCreator interface {
	CreateDeliver() Deliver
}

type Logistic struct {
	EmployeeAmount       int
	TotalTransportAmount int
	ListOfTransport      []Transport
	MapOfTransport       map[Transport]Order
}

func (t *Logistic) CreateDeliver() Deliver {
	fmt.Println("Creating a default vehicle for delivery...")
	return &Transport{}
}

type TrackLogistic struct {
	Logistic
}

func (t *TrackLogistic) CreateDeliver() Deliver {
	fmt.Println("Creating a truck for delivery...")
	return &Truck{
		TiresAmount: 12,
		Transport:   Transport{PassengerCapacity: 2, MaximumLoad: 1000, Fuel: 250, EnginePower: 400},
	}
}

type ShipLogistic struct {
	Logistic
}

func (s *ShipLogistic) CreateDeliver() Deliver {
	fmt.Println("Creating a ship for delivery...")
	return &Ship{
		SeaMileage: 500,
		Transport:  Transport{PassengerCapacity: 100, MaximumLoad: 2000, Fuel: 1500, EnginePower: 800},
	}
}

func main() {
	var creator DeliverCreator

	creator = &Logistic{}
	creator.CreateDeliver()

	creator = &ShipLogistic{}
	ship := creator.CreateDeliver()
	ship.Delivery("Port A", "Port B", 150)

	creator = &TrackLogistic{}
	truck := creator.CreateDeliver()
	truck.Delivery("Warehouse A", "Warehouse B", 500)
}
