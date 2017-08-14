package idy

import "testing"

func Test_fisherYatesShuffle(t *testing.T) {
	s1 := []uint64{1}
	fisherYatesShuffle(s1, 1)
	if s1[0] != 1 {
		t.Error("failed to shuffle slice with 1 element")
	}

	s2 := []uint64{1, 2, 3, 4, 5, 6}
	s3 := []uint64{1, 2, 3, 4, 5, 6}
	fisherYatesShuffle(s2, 1234)
	fisherYatesShuffle(s3, 1234)

	equal := true
	for i := 0; i < len(s2); i = i + 1 {
		if s2[i] != s3[i] {
			equal = false
		}
	}

	if !equal {
		t.Error("two shuffle not equal", s2, s3)
	}
}
