package main

import (
	"fmt"
	"io"
	"os"
)

const aA = 'a' - 'A'

func equ(l, r []byte) bool {
	ll, lr := len(l), len(r)
	if ll != lr {
		return false
	}
	for i := 0; i < ll; i++ {
		if l[i] != r[i] {
			return false
		}
	}
	return true
}

// compares two byte strings
func cmp(l, r []byte) int {
	ll, lr := len(l), len(r)
	len := ll
	if lr < len {
		len = lr
	}
	for i := 0; i < len; i++ {
		if l[i] < r[i] {
			return -1
		}
		if l[i] > r[i] {
			return 1
		}
	}
	if ll < lr {
		return -1
	}
	if lr < ll {
		return 1
	}
	return 0
}

type entry struct {
	s []byte
	n int
	// left and right children
	l, r int
}

// wordMap combines list and binary tree in place
type wordMap struct {
	d []entry
}

func newWordMap() *wordMap {
	return &wordMap{
		d: []entry{
			// root node
			entry{l: -1, r: -1},
		},
	}
}

// inc adds one to the count associated with the byte string
func (m *wordMap) inc(s []byte) {
	// starting at the root node ...
	n := 0

	for {
		// compare the new string
		d := cmp(s, m.d[n].s)
		if d == 0 {
			// if match, increment
			m.d[n].n++
			return
		} else if d < 0 {
			// else maybe descend the left branch
			nx := m.d[n].l
			if nx == -1 {
				// if there is no left branch, then we will make one
				m.d[n].l = len(m.d)
				break
			}
			n = nx
		} else {
			// same on the right
			nx := m.d[n].r
			if nx == -1 {
				m.d[n].r = len(m.d)
				break
			}
			n = nx
		}
	}

	// if not returned, then need a new node

	// s is a temporary buffer, so copy
	c := make([]byte, len(s))
	copy(c, s)

	// insert the new node, the parent will already be pointing to it
	m.d = append(m.d, entry{s: c, n: 1, l: -1, r: -1})
}

// inserts an entry into the right place in a sorted list
func insertIntoTopList(list []entry, toInsert entry) {
	// find the first entry the new one is bigger than
	for i, e := range list {
		if toInsert.n > e.n {
			// then move all subsequent entries down
			for i2 := len(list) - 1; i2 > i; i2-- {
				// starting at the end, so that we don't overwrite
				list[i2] = list[i2-1]
			}
			list[i] = toInsert
			return
		}
	}
	// if we haven't returned, then this must go in the last place
	list[len(list)] = toInsert
}

// topN returns a sorted list of the N most used words
func (m *wordMap) topN(count int) []entry {
	// topN can't be bigger than total words found
	if max := len(m.d) - 1; count > max {
		count = max
	}

	out := make([]entry, count)

	for _, e := range m.d {
		// see if this entry could be in the topN
		if e.n > out[count-1].n {
			// if it can, then insert in tht correct place
			insertIntoTopList(out, e)
		}
	}

	return out
}

// fatalf just a helper for crashing out
func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		fatalf("invalid call")
	}
	fn := os.Args[1]

	fi, err := os.Open(fn)
	if err != nil {
		fatalf("open error: %v", err)
	}
	defer fi.Close()

	// file buffer for reading through
	fb := make([]byte, 2<<12)
	// word buffer for copying each word, in standard form
	wb := make([]byte, 2<<8)
	// word slice for managing word buffer
	w := wb[0:0]
	// word map for storing results
	m := newWordMap()

	for {
		// read a chunk of the file
		n, err := fi.Read(fb)
		if err != nil {
			if n == 0 && err == io.EOF {
				break
			}
			// non EOF error is a problem
			fatalf("read error: %v", err)
		}
		// file slice for managing file buffer
		f := fb[0:n]
		for _, b := range f {
			if b >= 'a' && b <= 'z' {
				// lower case letters are part of this word
				w = append(w, b)
			} else if b >= 'A' && b <= 'Z' {
				// upper case letters are lowered and are then part of this word
				w = append(w, b+aA)
			} else {
				// anything else means the end of the word
				if len(w) > 0 {
					// if the word has any letters, then incremement the count
					m.inc(w)
					// and reset the word slice
					w = wb[0:0]
				}
			}
		}
	}

	for _, e := range m.topN(20) {
		fmt.Printf("% 7d %s\n", e.n, e.s)
	}
}
