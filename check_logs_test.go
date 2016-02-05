package main

import (
	"testing"
)

func TestAddNull(t *testing.T) {
	s := "7"
	if s = addNull(s); s != "07" {
		t.Fatalf("addNull doesn't work!")
	}
}
func TestHoursInterval(t *testing.T) {

	var tests = []struct {
		input []int
		want  []string
	}{
		{[]int{5, 20}, []string{"20", "19", "18", "17", "16", "15"}},
		{[]int{7, 5}, []string{"05", "04", "03", "02", "01", "00"}},
		{[]int{5, 0}, []string{"00"}},
	}

	for _, test := range tests {
		if got := hoursInterval(test.input[0], test.input[1]); len(got) != len(test.want) {
			t.Errorf(" hourseInterval(% q) = %v", test.input, got)
		}
	}
}
