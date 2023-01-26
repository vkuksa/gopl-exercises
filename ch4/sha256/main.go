// Exercises 4.1, 4.2 from gopl.io

package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)

const (
	SHA256 = iota
	SHA384
	SHA512
)

var encryption = flag.Int("e", SHA256, "0 for SHA256[default]\n1 for SHA384\n2 for SHA512")
var value = flag.String("v", "", "Input string to be encrypted")

func main() {
	flag.Parse()

	if len(*value) == 0 {
		fmt.Print("Provided empty string for encryption")
		os.Exit(1)
	}

	fmt.Println(*value)

	switch *encryption {
	case SHA384:
		fmt.Println(sha512.Sum384([]byte(*value)))
	case SHA512:
		fmt.Println(sha512.Sum512([]byte(*value)))
	default:
		fmt.Println(sha256.Sum256([]byte(*value)))
	}
}

// Counts the number of bits that are different between hashes
func bitDiff(c1, c2 [32]byte) uint8 {
	var result uint8 = 0

	for i := 0; i < 32; i++ {
		for b := 0; b < 8; b++ {
			if (c1[i] >> b) != (c2[i] >> b) {
				result++
			}
		}
	}

	return result
}
