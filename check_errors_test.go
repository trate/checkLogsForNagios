package main

import (
	"testing"
)

func TestAddNull (t *testing.T) {
	s := "7"
	if s = addNull(s); s != "07" {
		t.Fatalf("addNull doesn't work!")
	}
}
func TestHoursInterval (t *testing.T) {
	interval := []string{"20", "19", "18", "17", "16", "15",}
	interval2 := []string{"05", "04", "03", "02", "01", "00",}
	
	if len(interval) != len(hoursInterval(5, 20)) || len(interval2) != len(hoursInterval(7, 5)) {
		t.Fatalf("hoursInterval is not working correctly! %d, %d %v, %d, %d ", len(interval), len(hoursInterval(15,20)), hoursInterval(15,20), len(interval2), len(hoursInterval(7, 5)))
	} 	
}
