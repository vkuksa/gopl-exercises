// Exercise 2.2 from gopl.io

package convert

import (
	"fmt"
	"os"
	"strconv"
)

var (
	value int
)

func main() {
	if len(os.Args) == 4 {
		value, _ = strconv.Atoi(os.Args[3])
	}

	switch os.Args[1] {
	case "t":
		switch os.Args[2] {
		case "f":
			result := FtoC(Fahrenheit(value))
			fmt.Println(result)
		case "c":
			result := CtoF(Celsius(value))
			fmt.Println(result)
		default:
			fmt.Println("Expected f or c")
		}
	case "l":
		switch os.Args[2] {
		case "f":
			result := FtoM(Feet(value))
			fmt.Println(result)
		case "m":
			result := MtoF(Meter(value))
			fmt.Println(result)
		default:
			fmt.Println("Expected f or m")
		}
	case "w":
		switch os.Args[2] {
		case "p":
			result := PtoK(Pound(value))
			fmt.Println(result)
		case "k":
			result := KtoP(Kilogram(value))
			fmt.Println(result)
		default:
			fmt.Println("Expected f or c")
		}
	default:
		fmt.Println("Expected t, l, w")
		os.Exit(1)
	}
}
