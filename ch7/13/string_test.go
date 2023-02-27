package eval

import (
	"testing"
)

// !+String
func TestString(t *testing.T) {
	tests := []string{
		"sqrt(A / pi)",
		"pow(x, -3) + pow(y, 3.14)",
		"5 / 9 * (F - 32)",
		"5 / 9 * (F - 32)",
		"5 / 9 * (F - 32)",
	}
	for _, input := range tests {
		expr, err := Parse(input)
		if err != nil {
			t.Error(err) // parse error
			continue
		}
		expr, err = Parse(expr.String())
		if err != nil {
			t.Errorf("Parsing String() for %s, got %s\n, error %e", input, expr.String(), err)
		}
	}
}

//!-String
