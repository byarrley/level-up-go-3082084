package main

import (
	"testing"
)

type testCase struct {
	name  string
	input string
	want  bool
}

func TestIsBalanced(t *testing.T) {
	tcs := []testCase{
		{"balanced", "[1 + (2 + {3})]", true},
		{"No paren", "1+2", true},
		{"Extra left", "[[1 + 2]", false},
		{"Extra right", "(1 + 2))", false},
		{"Offset balanced", "[(1] + 2)", false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := isBalanced(tc.input)
			if got != tc.want {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}
