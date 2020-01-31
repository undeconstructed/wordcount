package main

import "testing"

func TestCmp(t *testing.T) {
	if -1 != cmp([]byte{1, 2, 2}, []byte{1, 2, 3}) {
		t.Fail()
	}
	if -1 != cmp([]byte{1, 2}, []byte{1, 2, 3}) {
		t.Fail()
	}
	if 1 != cmp([]byte{1, 2, 4}, []byte{1, 2, 3}) {
		t.Fail()
	}
	if 1 != cmp([]byte{1, 2, 3, 1}, []byte{1, 2, 3}) {
		t.Fail()
	}
	if 0 != cmp([]byte{1, 2, 3}, []byte{1, 2, 3}) {
		t.Fail()
	}
	if 0 != cmp([]byte{}, []byte{}) {
		t.Fail()
	}
	if 1 != cmp([]byte{0}, []byte{}) {
		t.Fail()
	}
}

func TestWordMap(t *testing.T) {
	m := newWordMap()
	m.inc([]byte{2})
	m.inc([]byte{1})
	m.inc([]byte{1})
	m.inc([]byte{1})
	m.inc([]byte{3, 4, 5})
	r := m.topN(100)
	if 3 != len(r) {
		t.Fail()
	}
	if 1 != r[0].s[0] {
		t.Fail()
	}
	if 3 != r[0].n {
		t.Fail()
	}
	if 1 != r[1].n {
		t.Fail()
	}
}
