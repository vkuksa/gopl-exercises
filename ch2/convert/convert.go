// Exercise 2.1 from gopl.io

package convert

import "fmt"

type Celsius float64
type Fahrenheit float64

type Meter float64
type Feet float64

type Pound float64
type Kilogram float64

func CtoF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

func FtoC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

func MtoF(m Meter) Feet {
	return Feet(m * 3.28)
}

func FtoM(f Feet) Meter {
	return Meter(f / 3.28)
}

func PtoK(p Pound) Kilogram {
	return Kilogram(p / 2.2046)
}

func KtoP(k Kilogram) Pound {
	return Pound(k * 2.2046)
}

func (c Celsius) String() string    { return fmt.Sprintf("%g°C", c) }
func (f Fahrenheit) String() string { return fmt.Sprintf("%g°F", f) }
func (m Meter) String() string      { return fmt.Sprintf("%g°C", m) }
func (f Feet) String() string       { return fmt.Sprintf("%g°F", f) }
func (p Pound) String() string      { return fmt.Sprintf("%g°C", p) }
func (k Kilogram) String() string   { return fmt.Sprintf("%g°F", k) }
