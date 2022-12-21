// Exercises 2.3 2.4 2.5 from gopl.io

package popcount

// pc[i] is the population count of i.
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func Calculate(x uint64) int {
	result := 0
	for i := 0; i < 8; i++ {
		result += (int(x) >> (i * 8))
	}

	return result
}

func CalculateByClearing(x uint64) int {
	n := 0
	for x != 0 {
		x = x & (x - 1) // clear rightmost non-zero bit
		n++
	}
	return n
}

func CalculateByShifting(x uint64) int {
	n := 0
	for i := uint(0); i < 64; i++ {
		if x&(1<<i) != 0 {
			n++
		}
	}
	return n
}

//!-
