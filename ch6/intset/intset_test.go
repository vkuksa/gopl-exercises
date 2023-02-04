package intset

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	var testWords = []uint{4398046511618, 0, 65536}

	var x IntSet
	x.words = make([]uint, 3)
	copy(x.words, testWords)

	if fmt.Sprintf(x.String()) != "{1 9 42 144}" {
		t.Fatalf("String")
	}
}

func TestAddRemoveHasElems(t *testing.T) {
	var x IntSet
	x.Add(1)
	x.Add(144)
	if fmt.Sprintf(x.String()) != "{1 144}" {
		t.Fatalf("Add")
	}

	x.AddAll(9, 515)
	if fmt.Sprintf(x.String()) != "{1 9 144 515}" {
		t.Fatalf("AddAll")
	}

	x.Remove(144)
	if fmt.Sprintf(x.String()) != "{1 9 515}" {
		t.Fatalf("Remove 144")
	}

	x.Remove(515)
	if fmt.Sprintf(x.String()) != "{1 9}" {
		t.Fatalf("Remove 515")
	}

	elems := x.Elems()
	if len(elems) != 2 || elems[0] != 1 || elems[1] != 9 {
		t.Fatalf("Elems")
	}

	if !x.Has(1) || !x.Has(9) || x.Has(515) || x.Has(144) || x.Has(123) {
		t.Fatalf("Has")
	}
}

func TestLenClearCopy(t *testing.T) {
	var x IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)

	if x.Len() != 3 {
		t.Errorf("Len")
	}

	y := x.Copy()
	if y.Len() != 3 && x.Len() != 3 {
		t.Errorf("Copy")
	}

	x.Clear()
	if x.Len() != 0 && y.Len() == 0 {
		t.Errorf("Clear")
	}
}

func TestBitwiseOperations(t *testing.T) {
	{
		var x, y IntSet
		x.Add(1)
		x.Add(144)
		x.Add(9)
		y.Add(9)
		y.Add(42)

		union := (*IntSet).UnionWith
		union(&y, &x)
		if fmt.Sprintf(y.String()) != "{1 9 42 144}" {
			t.Errorf("UnionWith y")
		}
		if fmt.Sprintf(x.String()) != "{1 9 144}" {
			t.Errorf("UnionWith x")
		}
	}
	{
		var x, y IntSet
		x.Add(1)
		x.Add(144)
		x.Add(9)
		y.Add(9)
		y.Add(42)

		intersect := (*IntSet).IntersectWith
		intersect(&y, &x)
		if fmt.Sprintf(y.String()) != "{9}" {
			t.Errorf("IntersectWith y")
		}
		if fmt.Sprintf(x.String()) != "{1 9 144}" {
			t.Errorf("IntersectWith x")
		}
	}
	{
		var x, y IntSet
		x.Add(1)
		x.Add(144)
		x.Add(9)
		y.Add(9)
		y.Add(42)

		difference := (*IntSet).DifferenceWith
		difference(&y, &x)
		if fmt.Sprintf(y.String()) != "{42}" {
			t.Errorf("DifferenceWith y")
		}
		if fmt.Sprintf(x.String()) != "{1 9 144}" {
			t.Errorf("DifferenceWith x")
		}
	}
	{
		var x, y IntSet
		x.Add(1)
		x.Add(144)
		x.Add(9)
		y.Add(9)
		y.Add(42)

		symdiff := (*IntSet).SymmetricDifference
		symdiff(&y, &x)
		if fmt.Sprintf(y.String()) != "{1 42 144}" {
			t.Errorf("SymmetricDifference y")
		}
		if fmt.Sprintf(x.String()) != "{1 9 144}" {
			t.Errorf("SymmetricDifference x")
		}
	}
}
