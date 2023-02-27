package eval

import (
	"bytes"
	"fmt"
	"strconv"
)

//!+String

func (v Var) String() string {
	return string(v)
}

func (l literal) String() string {
	return strconv.FormatFloat(float64(l), 'f', -1, 64)
}

func (u unary) String() string {
	return fmt.Sprintf("(%c%s)", u.op, u.x.String())
}

func (b binary) String() string {
	return fmt.Sprintf("(%s %c %s)", b.x.String(), b.op, b.y.String())
}

func (c call) String() string {
	var buf = new(bytes.Buffer)
	fmt.Fprintf(buf, "%s(", c.fn)
	for i, arg := range c.args {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	buf.WriteByte(')')

	return buf.String()
}

//!-String
