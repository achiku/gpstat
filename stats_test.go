package main

import "testing"

func TestAppend(t *testing.T) {
	s := NewStatsList()
	s.Append()
	s.Append()
	t.Logf("%v", s.Data)
}
