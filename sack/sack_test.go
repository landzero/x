package sack

import "testing"

func TestSack_ParentGetParent(t *testing.T) {
	p := Sack{"a": "b"}
	s := Sack{}
	s.SetParent(p)
	if s.Parent().Value("a") != "b" {
		t.Errorf("bad parent")
	}
	var n Sack
	n.Parent()     // should not PANIC
	n.SetParent(p) // should not PANIC
}

func TestSack_SetValue(t *testing.T) {
	z := Sack{"0": "1"}
	p := Sack{"a": "b"}
	p.SetParent(z)
	s := Sack{"c": "d"}
	s.SetParent(p)
	if s.Value("c") != "d" {
		t.Errorf("bad key c")
	}
	if s.Value("a") != "b" {
		t.Errorf("bad key a")
	}
	if s.Value("0") != "1" {
		t.Errorf("bad key 0")
	}
	var n Sack
	n.Set("e", "f") // should not PANIC
	n.Value("h")    // should not PANIC
}
