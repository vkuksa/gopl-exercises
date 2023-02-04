package intset

import (
	"bytes"
	"fmt"
	"math/bits"
)

const wordSize = bits.UintSize

type IntSet struct {
	words []uint
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *IntSet) shrinkToFit() {
	for i := len(s.words) - 1; i > 0; i-- {
		for j := 0; j < wordSize; j++ {
			bit := s.words[i] & (1 << uint(j))
			if bit != 0 {
				return
			}
		}
		s.words = s.words[:len(s.words)-1]
	}
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/wordSize, uint(x%wordSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/wordSize, uint(x%wordSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) AddAll(values ...int) {
	for _, value := range values {
		s.Add(value)
	}
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) IntersectWith(t *IntSet) {
	minLen := min(len(t.words), len(s.words))

	for i := 0; i < minLen; i++ {
		s.words[i] &= t.words[i]
	}

	s.words = s.words[:minLen]
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	minLen := min(len(t.words), len(s.words))

	for i := 0; i < minLen; i++ {
		s.words[i] &= ^t.words[i]
	}

	s.words = s.words[:minLen]
}

func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// remove x from the set
func (s *IntSet) Remove(x int) {
	word, bit := x/wordSize, uint(x%wordSize)
	s.words[word] &= ^(1 << bit)
	s.shrinkToFit()
}

// return the number of elements
func (s *IntSet) Len() int {
	return len(s.words)
}

// remove all elements from the set
func (s *IntSet) Clear() {
	s.words = make([]uint, 0)
}

// return a copy of the set
func (src *IntSet) Copy() (cp *IntSet) {
	cp = &IntSet{}
	cp.words = make([]uint, len(src.words))
	copy(src.words, cp.words)
	return cp
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", wordSize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func (s *IntSet) Elems() []int {
	elems := make([]int, 0)

	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordSize; j++ {
			if word&(1<<uint(j)) != 0 {
				elems = append(elems, wordSize*i+j)
			}
		}
	}

	return elems
}
