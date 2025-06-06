package main

import "fmt"

type Builder interface {
	SetSkin()
	SetEyes()
	SetLegs()
	SetArms()
	SetMouth()
	SetHair()
	SetHeight()
	GetPerson() Person
}

func getBuilder(s string) Builder {
	switch s {
	case "asia":
		return NewAsianBuilder()
	case "europe":
		return NewEuropeanBuilder()
	default:
		fmt.Print("Unknown Builder")
		return nil
	}
}

type AsianBuilder struct {
	skin   string
	eyes   string
	legs   int
	arms   int
	mouth  string
	hair   string
	height float32
}

func NewAsianBuilder() *AsianBuilder {
	return &AsianBuilder{}
}

func (a *AsianBuilder) SetSkin() {
	a.skin = "Yellow"
}

func (a *AsianBuilder) SetEyes() {
	a.eyes = "Black"
}

func (a *AsianBuilder) SetLegs() {
	a.legs = 2
}

func (a *AsianBuilder) SetArms() {
	a.arms = 2
}

func (a *AsianBuilder) SetMouth() {
	a.mouth = "Small"
}

func (a *AsianBuilder) SetHair() {
	a.hair = "Black"
}

func (a *AsianBuilder) SetHeight() {
	a.height = 1.6
}

func (a AsianBuilder) GetPerson() Person {
	return Person{
		skin:   a.skin,
		eyes:   a.eyes,
		legs:   a.legs,
		arms:   a.arms,
		mouth:  a.mouth,
		hair:   a.hair,
		height: a.height,
	}
}

type EuropeanBuilder struct {
	skin   string
	eyes   string
	legs   int
	arms   int
	mouth  string
	hair   string
	height float32
}

func NewEuropeanBuilder() *EuropeanBuilder {
	return &EuropeanBuilder{}
}

func (a *EuropeanBuilder) SetSkin() {
	a.skin = "White"
}

func (a *EuropeanBuilder) SetEyes() {
	a.eyes = "Brown"
}

func (a *EuropeanBuilder) SetLegs() {
	a.legs = 2
}

func (a *EuropeanBuilder) SetArms() {
	a.arms = 2
}

func (a *EuropeanBuilder) SetMouth() {
	a.mouth = "Medium"
}

func (a *EuropeanBuilder) SetHair() {
	a.hair = "Black"
}

func (a *EuropeanBuilder) SetHeight() {
	a.height = 1.75
}

func (a EuropeanBuilder) GetPerson() Person {
	return Person{
		skin:   a.skin,
		eyes:   a.eyes,
		legs:   a.legs,
		arms:   a.arms,
		mouth:  a.mouth,
		hair:   a.hair,
		height: a.height,
	}
}
