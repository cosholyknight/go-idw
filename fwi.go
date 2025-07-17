package idw

import (
	"errors"
	"math"
)

func FFMC(temp, rh, wind, rain, ffmcPrev float64) float64 {
	rh = math.Min(rh, 100.0)
	mo := 147.2 * (101.0 - ffmcPrev) / (59.5 + ffmcPrev)

	if rain > 0.5 {
		rf := rain - 0.5
		var mr float64
		if mo <= 150.0 {
			mr = mo + 42.5*rf*math.Exp(-100.0/(251.0-mo))*(1.0-math.Exp(-6.93/rf))
		} else {
			mr = mo + 42.5*rf*math.Exp(-100.0/(251.0-mo))*(1.0-math.Exp(-6.93/rf)) +
				0.0015*math.Pow(mo-150.0, 2)*math.Sqrt(rf)
		}
		if mr > 250.0 {
			mr = 250.0
		}
		mo = mr
	}

	ed := 0.942*math.Pow(rh, 0.679) + 11.0*math.Exp((rh-100.0)/10.0) + 0.18*(21.1-temp)*(1.0-math.Exp(-0.115*rh))

	var m float64
	if mo > ed {
		ko := 0.424*(1.0-math.Pow(rh/100.0, 1.7)) + 0.0694*math.Sqrt(wind)*(1.0-math.Pow(rh/100.0, 8))
		kd := ko * 0.581 * math.Exp(0.0365*temp)
		m = ed + (mo-ed)*math.Pow(10.0, -kd)
	} else {
		ew := 0.618*math.Pow(rh, 0.753) + 10.0*math.Exp((rh-100.0)/10.0) + 0.18*(21.1-temp)*(1.0-math.Exp(-0.115*rh))
		if mo < ew {
			k1 := 0.424*(1.0-math.Pow((100.0-rh)/100.0, 1.7)) + 0.0694*math.Sqrt(wind)*(1.0-math.Pow((100.0-rh)/100.0, 8))
			kw := k1 * 0.581 * math.Exp(0.0365*temp)
			m = ew - (ew-mo)*math.Pow(10.0, -kw)
		} else {
			m = mo
		}
	}
	return 59.5 * (250.0 - m) / (147.2 + m)
}

func DayLength(latitude float64, month int) (float64, error) {
	if month < 1 || month > 12 {
		return 0, errors.New("invalid month")
	}
	var dayLengths []float64
	switch {
	case latitude > 33 && latitude <= 90:
		dayLengths = []float64{6.5, 7.5, 9.0, 12.8, 13.9, 13.9, 12.4, 10.9, 9.4, 8.0, 7.0, 6.0}
	case latitude > 0 && latitude <= 33:
		dayLengths = []float64{7.9, 8.4, 8.9, 9.5, 9.9, 10.2, 10.1, 9.7, 9.1, 8.6, 8.1, 7.8}
	case latitude > -30 && latitude <= 0:
		dayLengths = []float64{10.1, 9.6, 9.1, 8.5, 8.1, 7.8, 7.9, 8.3, 8.9, 9.4, 9.9, 10.2}
	case latitude >= -90 && latitude <= -30:
		dayLengths = []float64{11.5, 10.5, 9.2, 7.9, 6.8, 6.2, 6.5, 7.4, 8.7, 10.0, 11.2, 11.8}
	default:
		return 0, errors.New("invalid latitude")
	}
	return dayLengths[month-1], nil
}

func DryingFactor(lat float64, month int) float64 {
	LfN := []float64{-1.6, -1.6, -1.6, 0.9, 3.8, 5.8, 6.4, 5.0, 2.4, 0.4, -1.6, -1.6}
	LfS := []float64{6.4, 5.0, 2.4, 0.4, -1.6, -1.6, -1.6, -1.6, -1.6, 0.9, 3.8, 5.8}
	if lat > 0 {
		return LfN[month-1]
	}
	return LfS[month-1]
}

func DMC(temp, rh, rain, dmcPrev, lat float64, month int) float64 {
	rh = math.Min(100.0, rh)

	if rain > 1.5 {
		re := 0.92*rain - 1.27
		mo := 20.0 + math.Exp(5.6348-dmcPrev/43.43)

		var b float64
		switch {
		case dmcPrev <= 33.0:
			b = 100.0 / (0.5 + 0.3*dmcPrev)
		case dmcPrev <= 65.0:
			b = 14.0 - 1.3*math.Log(dmcPrev)
		default:
			b = 6.2*math.Log(dmcPrev) - 17.2
		}

		mr := mo + 1000.0*re/(48.77+b*re)
		pr := 244.72 - 43.43*math.Log(mr-20.0)
		if pr > 0 {
			dmcPrev = pr
		} else {
			dmcPrev = 0.0
		}
	}

	if temp > -1.1 {
		d1, _ := DayLength(lat, month)
		k := 1.894 * (temp + 1.1) * (100.0 - rh) * d1 * 1e-6
		dmcPrev += 100 * k
	}

	return dmcPrev
}

func DC(temp, rain, dcPrev, lat float64, month int) float64 {
	if rain > 2.8 {
		rd := 0.83*rain - 1.27
		qo := 800.0 * math.Exp(-dcPrev/400.0)
		qr := qo + 3.937*rd
		dr := 400.0 * math.Log(800.0/qr)
		if dr > 0 {
			dcPrev = dr
		} else {
			dcPrev = 0.0
		}
	}

	lf := DryingFactor(lat, month)
	v := 0.36*(temp+2.8) + lf
	if temp <= -2.8 {
		v = lf
	}
	if v < 0.0 {
		v = 0.0
	}

	return dcPrev + 0.5*v
}

func ISI(wind, ffmc float64) float64 {
	fW := math.Exp(0.05039 * wind)
	m := 147.2 * (101.0 - ffmc) / (59.5 + ffmc)
	fF := 91.9 * math.Exp(-0.1386*m) * (1.0 + math.Pow(m, 5.31)/49300000.0)
	return 0.208 * fW * fF
}

func BUI(dmc, dc float64) float64 {
	if dmc <= 0.4*dc {
		return 0.8 * dmc * dc / (dmc + 0.4*dc)
	}
	return dmc - (1.0-0.8*dc/(dmc+0.4*dc))*(0.92+math.Pow(0.0114*dmc, 1.7))
}

func FWI(isi, bui float64) float64 {
	var fD float64
	if bui <= 80.0 {
		fD = 0.626*math.Pow(bui, 0.809) + 2.0
	} else {
		fD = 1000.0 / (25.0 + 108.64*math.Exp(-0.023*bui))
	}
	b := 0.1 * isi * fD
	if b > 1.0 {
		return math.Exp(2.72 * math.Pow(0.434*math.Log(b), 0.647))
	}
	return b
}

func CalcFWI(month int, temp, rh, wind, rain, ffmcPrev, dmcPrev, dcPrev, lat float64) (float64, error) {
	ffmc := FFMC(temp, rh, wind, rain, ffmcPrev)
	dmc := DMC(temp, rh, rain, dmcPrev, lat, month)
	dc := DC(temp, rain, dcPrev, lat, month)
	isi := ISI(wind, ffmc)
	bui := BUI(dmc, dc)
	fwi := FWI(isi, bui)
	return fwi, nil
}
