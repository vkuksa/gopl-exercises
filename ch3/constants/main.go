// Exercise 3.13 from gopl.io

package main

import (
	"fmt"
)

// const (
// 	Byte = 1 << (10 * iota)
// 	KiloByte
// 	MegaByte
// 	GigaByte
// 	TerraByte
// 	PetaByte
// 	ExaByte
// 	ZettaByte
// 	YottaByte
// )

const (
	Step      = 1000
	Byte      = 1
	KiloByte  = Byte * Step
	MegaByte  = KiloByte * Step
	GigaByte  = MegaByte * Step
	TerraByte = GigaByte * Step
	PetaByte  = TerraByte * Step
	ExaByte   = PetaByte * Step
	ZettaByte = ExaByte * Step
	YottaByte = ZettaByte * Step
)

func main() {
	fmt.Println(Byte)
	fmt.Println(KiloByte)
	fmt.Println(MegaByte)
	fmt.Println(GigaByte)
	fmt.Println(TerraByte)
	fmt.Println(PetaByte)
	fmt.Println(ExaByte)
	// fmt.Println(ZettaByte)
	// fmt.Println(YottaByte)

}

//!-
