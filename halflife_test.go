package main

import (
	"testing"
	"time"
)

func TestHalves(t *testing.T) {

	tables := []struct {
		initial int
		halves int
	} {
		{100, 7},
		{200, 8},
		{5, 2},
	}

	for _, table := range tables {
		caffHalves := Halves(table.initial)
		if caffHalves != table.halves {
			t.Errorf("Halvgorithm of %d was incorrect, got: %d, want: %d.",
				table.initial, caffHalves, table.halves)
		}
	}
}

func TestMetabolize(t *testing.T) {
	dosage := 100
	bufSize := Halves(dosage)
	dose := Dose { Name: "Chai Tea Latte", Dosage: dosage, Time: time.Now().Format(time.RFC3339) }
	caffeine := Metabolizer { Name: "caffeine", Onset: time.Hour*1, Halflife: (time.Hour*5 + time.Minute*42) }

	channel := make(chan Remnant, bufSize)
	go Metabolize(dose, caffeine, channel)
	for remnant := range channel {
		if remnant.Amount >= dosage {
			t.Errorf("Remnant calculation of %d was incorrect, got %d",
				dosage, remnant.Amount)
		}
	}
}
