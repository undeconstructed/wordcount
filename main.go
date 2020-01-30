package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

const aA = 'a' - 'A'

type entry struct {
	s []byte
	n int
}

type wordMap struct {
	d []entry
}

func (m *wordMap) inc(s []byte) {

	// XXX horrible slow lookup
check:
	for n, e := range m.d {
		if len(e.s) == len(s) {
			for i := 0; i < len(s); i++ {
				if e.s[i] != s[i] {
					continue check
				}
			}
			m.d[n].n++
			return
		}
	}

	// s is a temporary buffer, so copy
	c := make([]byte, len(s))
	copy(c, s)
	m.d = append(m.d, entry{
		s: c,
		n: 1,
	})
}

func (m *wordMap) topTwenty() []entry {
	// XXX library function to remove
	sort.Slice(m.d, func(i int, j int) bool {
		return m.d[i].n > m.d[j].n
	})

	n := 20
	if n > len(m.d) {
		n = len(m.d)
	}

	return m.d[0:n]
}

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

	// file buffer for readin through
	fb := make([]byte, 2<<12)
	// word buffer for copying each word, in standard form
	wb := make([]byte, 2<<8)
	// word slice for managing word buffer
	w := wb[0:0]
	// word map for storing results
	m := &wordMap{}

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

	for _, e := range m.topTwenty() {
		fmt.Printf("% 7d %s\n", e.n, e.s)
	}
}
