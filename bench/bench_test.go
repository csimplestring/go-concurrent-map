package main

import "testing"

func TestRandKey(t *testing.T) {
	s1 := randKey(-1)
	s2 := randKey(-1)
	s3 := randKey(-1)
	s4 := randKey(-1)

	strings := []string{s1, s2, s3, s4}
	for _, s := range strings {
		t.Log(s)
	}
}
