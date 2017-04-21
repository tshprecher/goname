package main

import (
	"testing"
)

func TestName(t *testing.T) {
	tests := []struct {
		in string
		ex string
	}{
		{"ONE_TWO_THREE", "OneTwoThree"},
		{"one_two_three", "oneTwoThree"},
		{"oneTwoThree", "oneTwoThree"},
		{"ALLCAPS", "Allcaps"},
		{"_", "_"},
	}

	for _, tst := range tests {
		ex := rename(tst.in)
		if ex != tst.ex {
			t.Errorf("expected %s, received %s", tst.ex, ex)
		}
	}
}
