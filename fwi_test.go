package idw

import (
	"math"
	"testing"
)

func nearlyEqual(a, b float64) bool {
	return math.Abs(a-b) < 5e-5
}

type FwiInput struct {
	Month       int
	Day         int
	Temperature float64
	Humidity    float64
	Wind        float64
	Rain        float64
}

type FwiOutput struct {
	FFMC float64
	DMC  float64
	DC   float64
	ISI  float64
	BUI  float64
	FWI  float64
}

func TestFWICalculation(t *testing.T) {
	inputs := []FwiInput{
		{4, 13, 17, 42, 25, 0},
		{4, 14, 20, 21, 25, 2.4},
		{4, 15, 8.5, 40, 17, 1},
	}

	expected := []FwiOutput{
		{85, 6, 15, 0.0, 0.0, 0.0},
		{85.93646, 8.21537, 19.55400, 8.459204, 8.192343, 8.0247},
		{86.24764, 10.41108, 23.568, 8.83757, 10.36413, 9.28387},
	}

	prevFFMC := expected[0].FFMC
	prevDMC := expected[0].DMC
	prevDC := expected[0].DC

	// test on middle date: 14/4
	in := inputs[1]
	exp := expected[1]

	ffmc := FFMC(in.Temperature, in.Humidity, in.Wind, in.Rain, prevFFMC)
	dmc := DMC(in.Temperature, in.Humidity, in.Rain, prevDMC, 45.98, in.Month)
	dc := DC(in.Temperature, in.Rain, prevDC, 45.98, in.Month)
	isi := ISI(in.Wind, ffmc)
	bui := BUI(dmc, dc)
	fwi := FWI(isi, bui)

	if !nearlyEqual(ffmc, exp.FFMC) {
		t.Errorf("[Day %d] FFMC: got %f, expected %f", in.Day, ffmc, exp.FFMC)
	}
	if !nearlyEqual(dmc, exp.DMC) {
		t.Errorf("[Day %d] DMC: got %f, expected %f", in.Day, dmc, exp.DMC)
	}
	if !nearlyEqual(dc, exp.DC) {
		t.Errorf("[Day %d] DC: got %f, expected %f", in.Day, dc, exp.DC)
	}
	if !nearlyEqual(isi, exp.ISI) {
		t.Errorf("[Day %d] ISI: got %f, expected %f", in.Day, isi, exp.ISI)
	}
	if !nearlyEqual(bui, exp.BUI) {
		t.Errorf("[Day %d] BUI: got %f, expected %f", in.Day, bui, exp.BUI)
	}
	if !nearlyEqual(fwi, exp.FWI) {
		t.Errorf("[Day %d] FWI: got %f, expected %f", in.Day, fwi, exp.FWI)
	}
}
