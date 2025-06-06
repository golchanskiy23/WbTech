package main

type Director struct {
	builder Builder
}

func NewDirector(builder Builder) *Director {
	return &Director{builder}
}

func (d *Director) CreatePerson() Person {
	d.builder.SetArms()
	d.builder.SetEyes()
	d.builder.SetHair()
	d.builder.SetMouth()
	d.builder.SetHeight()
	d.builder.SetLegs()
	d.builder.SetSkin()
	return d.builder.GetPerson()
}
