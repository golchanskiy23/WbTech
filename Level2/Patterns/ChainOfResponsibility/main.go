package main

import (
	"fmt"
)

type Department interface {
	execute(*Request)
	setNext(Department)
}

type Illness = string

var (
	CAUGHT = "caught"
	FLU    = "flu"
	SHIZ   = "shiza"
)

type Patient struct {
	name, surname, age                                             string
	registrationDone, doctorCheckUpDone, medicineDone, paymentDone bool
	illness                                                        Illness
}

type Request struct {
	patient *Patient
}

type Reception struct {
	d               Department
	famousIllnesses map[Illness]struct{}
}

func (r *Reception) MapOfIllnessesInit() {
	r.famousIllnesses[CAUGHT] = struct{}{}
	r.famousIllnesses[FLU] = struct{}{}
	r.famousIllnesses[SHIZ] = struct{}{}
}

func (r *Reception) execute(req *Request) {
	if req.patient.illness == "" {
		fmt.Println("Person already healthy")
		return
	}
	if _, ok := r.famousIllnesses[req.patient.illness]; !ok {
		fmt.Println("Unfamous illness")
		return
	}
	fmt.Println("Reception stage accepted")
	req.patient.registrationDone = true
	r.d.execute(req)
}

func (r *Reception) setNext(department Department) {
	r.d = department
}

type Doctor struct {
	d Department
}

func (r *Doctor) execute(req *Request) {
	fmt.Println("Doctor stage accepted")
	req.patient.doctorCheckUpDone = true
	r.d.execute(req)
}

func (r *Doctor) setNext(department Department) {
	r.d = department
}

type MedicalProcedure struct {
	d Department
}

func (r *MedicalProcedure) execute(req *Request) {
	fmt.Println("Med procedure stage accepted")
	req.patient.medicineDone = true
	r.d.execute(req)
}

func (r *MedicalProcedure) setNext(department Department) {
	r.d = department
}

type Cashier struct {
	d Department
}

func (r *Cashier) execute(req *Request) {
	fmt.Println("Cashier stage accepted")
	req.patient.paymentDone = true
}

func (r *Cashier) setNext(department Department) {
	r.d = department
}

func main() {
	reception := &Reception{famousIllnesses: make(map[Illness]struct{})}
	reception.MapOfIllnessesInit()
	doctor := &Doctor{}
	procedure := &MedicalProcedure{}
	cashier := &Cashier{}

	reception.setNext(doctor)
	doctor.setNext(procedure)
	procedure.setNext(cashier)

	patients := []Patient{
		{name: "Max", surname: "Golchanskiy", illness: SHIZ},
		{name: "Tom", surname: "Peterson", illness: CAUGHT},
		{name: "Brad", surname: "Steel", illness: FLU},
	}
	for i, patient := range patients {
		currRequest := &Request{patient: &patient}
		reception.execute(currRequest)
		if patient.registrationDone && patient.doctorCheckUpDone && patient.medicineDone && patient.paymentDone {
			fmt.Printf("Patient number %d is healthy and payed for all services\n\n", i+1)
		} else {
			fmt.Printf("Something wrong happened with client number: %d!\n\n", i+1)
		}
	}
}
