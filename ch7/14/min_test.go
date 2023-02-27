package eval

import (
	"fmt"
	"testing"
)

// !+Eval
func TestMin(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"min(1, 2, 3)", Env{}, "1"},
		{"min(1, 2, 3, 4, -5, x)", Env{"x": -6}, "-6"},
		{"5 / 10 * min(6, 10, 15, 2, F)", Env{"F": 5}, "1"},
	}
	for _, test := range tests {
		fmt.Printf("\n%s\n", test.expr)

		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err) // parse error
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			t.Errorf("%s.Eval() in %v = %q, want %q\n",
				test.expr, test.env, got, test.want)
		}
	}
}

//!-Eval
