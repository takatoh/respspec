package response

import (
	"math"

	"github.com/takatoh/sdof/directintegration"
	"github.com/takatoh/seismicwave"
)

type Response struct {
	Period float64
	Sa     float64
	Sv     float64
	Sd     float64
}

func NewResponse(period, sa, sv, sd float64) *Response {
	p := new(Response)
	p.Period = period
	p.Sa = sa
	p.Sv = sv
	p.Sd = sd
	return p
}

func Spectrum(wave *seismicwave.Wave, period []float64, h float64) []*Response {
	var am, vm, dm float64

	spectrum := make([]*Response, 0)
	dt := wave.Dt / 10.0
	nperiod := len(period)
	z := interpolate(wave.Data, 10)

	for j := 0; j < nperiod; j++ {
		if math.Abs(period[j]) < 0.01 {
			am = 0.0
			for i := 1; i < len(z); i++ {
				if math.Abs(z[i]) > am {
					am = math.Abs(z[i])
				}
			}
			spectrum = append(spectrum, NewResponse(period[j], am, 0.0, 0.0))
		} else {
			am, vm, dm = WilsonTheta(z, dt, period[j], h)
			spectrum = append(spectrum, NewResponse(period[j], am, vm, dm))
		}
	}

	return spectrum
}

func Spectrum2(wave *seismicwave.Wave, period []float64, h float64) []*Response {
	var am, vm, dm float64

	spectrum := make([]*Response, 0)
	dt := wave.Dt / 10.0
	nperiod := len(period)
	z := interpolate(wave.Data, 10)
	n := len(z)

	for j := 0; j < nperiod; j++ {
		if math.Abs(period[j]) < 0.01 {
			am = 0.0
			for i := 1; i < len(z); i++ {
				if math.Abs(z[i]) > am {
					am = math.Abs(z[i])
				}
			}
			spectrum = append(spectrum, NewResponse(period[j], am, 0.0, 0.0))
		} else {
			w := 2.0 * math.Pi / period[j]
			acc, vel, dis := directintegration.WilsonTheta(h, w, dt, n, z)
			am = absMax(acc)
			vm = absMax(vel)
			dm = absMax(dis)
			spectrum = append(spectrum, NewResponse(period[j], am, vm, dm))
		}
	}

	return spectrum
}

func absMax(z []float64) float64 {
	zm := math.Abs(z[0])
	n := len(z)
	for i := 1; i < n; i++ {
		za := math.Abs(z[i])
		if za > zm {
			zm = za
		}
	}
	return zm
}

func interpolate(zin []float64, ndiv int) []float64 {
	var k int
	var zinc float64
	nin := len(zin)
	z := make([]float64, 0)
	k = 0
	z = append(z, 0.0)
	for i := 0; i < nin; i++ {
		if i == 0 {
			zinc = zin[i] / float64(ndiv)
		} else {
			zinc = (zin[i] - zin[i-1]) / float64(ndiv)
		}
		for j := 0; j < ndiv; j++ {
			z = append(z, z[k]+zinc)
			k++
		}
	}

	return z
}

// Wilson-theta method.
func WilsonTheta(z []float64, dt, period, h float64) (float64, float64, float64) {
	theta := 1.4

	tdt := theta * dt
	omega := 2.0 * math.Pi / period
	k := omega * omega
	c := 2.0 * h * omega

	// Constants for Willson-theta method.
	a1 := 1.0 + tdt*c/2.0 + k*tdt*tdt/6.0
	a2 := c + k*tdt
	a3 := tdt*c/2.0 + k/3.0*tdt*tdt

	// Set initial values.
	acc := 0.0
	vel := 0.0
	dis := 0.0
	am := 0.0
	vm := 0.0
	dm := 0.0

	for i := 1; i < len(z)-1; i++ {
		f := (theta-1.0)*z[i] - theta*z[i+1]
		ath := (f - k*dis - a2*vel - a3*acc) / a1
		acd := ((theta-1.0)*acc + ath) / theta
		dis = dis + dt*vel + acc*dt*dt/3.0 + acd*dt*dt/6.0
		vel = vel + (acc+acd)*dt/2.0
		acc = acd
		abz := math.Abs(acc + z[i+1])
		if abz > am {
			am = abz
		}
		if math.Abs(vel) > vm {
			vm = math.Abs(vel)
		}
		if math.Abs(dis) > dm {
			dm = math.Abs(dis)
		}
	}

	return am, vm, dm
}

func DefaultPeriod() []float64 {
	return []float64{
		0.02,
		0.02042,
		0.0208489,
		0.0212868,
		0.0217338,
		0.0221903,
		0.0226563,
		0.0231322,
		0.023618,
		0.024114,
		0.0246205,
		0.0251376,
		0.0256655,
		0.0262045,
		0.0267549,
		0.0273168,
		0.0278905,
		0.0284763,
		0.0290743,
		0.029685,
		0.0303084,
		0.0309449,
		0.0315949,
		0.0322584,
		0.0329359,
		0.0336276,
		0.0343339,
		0.035055,
		0.0357912,
		0.0365429,
		0.0373104,
		0.038094,
		0.038894,
		0.0397109,
		0.0405449,
		0.0413964,
		0.0422658,
		0.0431535,
		0.0440598,
		0.0449852,
		0.04593,
		0.0468946,
		0.0478795,
		0.0488851,
		0.0499118,
		0.05096,
		0.0520303,
		0.053123,
		0.0542387,
		0.0553779,
		0.0565409,
		0.0577284,
		0.0589408,
		0.0601787,
		0.0614426,
		0.062733,
		0.0640505,
		0.0653957,
		0.0667692,
		0.0681715,
		0.0696032,
		0.071065,
		0.0725576,
		0.0740814,
		0.0756373,
		0.0772258,
		0.0788477,
		0.0805037,
		0.0821945,
		0.0839207,
		0.0856832,
		0.0874828,
		0.0893201,
		0.091196,
		0.0931113,
		0.0950669,
		0.0970635,
		0.099102,
		0.101183,
		0.103308,
		0.105478,
		0.107693,
		0.109955,
		0.112264,
		0.114622,
		0.11703,
		0.119487,
		0.121997,
		0.124559,
		0.127175,
		0.129846,
		0.132573,
		0.135357,
		0.1382,
		0.141103,
		0.144066,
		0.147092,
		0.150181,
		0.153335,
		0.156556,
		0.159844,
		0.163201,
		0.166628,
		0.170128,
		0.173701,
		0.177349,
		0.181074,
		0.184877,
		0.188759,
		0.192724,
		0.196771,
		0.200904,
		0.205123,
		0.209432,
		0.21383,
		0.218321,
		0.222906,
		0.227588,
		0.232367,
		0.237248,
		0.24223,
		0.247318,
		0.252512,
		0.257815,
		0.26323,
		0.268758,
		0.274403,
		0.280166,
		0.28605,
		0.292058,
		0.298191,
		0.304454,
		0.310848,
		0.317377,
		0.324042,
		0.330848,
		0.337797,
		0.344891,
		0.352134,
		0.35953,
		0.367081,
		0.37479,
		0.382662,
		0.390698,
		0.398904,
		0.407282,
		0.415836,
		0.424569,
		0.433486,
		0.44259,
		0.451885,
		0.461376,
		0.471066,
		0.480959,
		0.49106,
		0.501374,
		0.511904,
		0.522655,
		0.533632,
		0.544839,
		0.556282,
		0.567965,
		0.579893,
		0.592072,
		0.604507,
		0.617203,
		0.630166,
		0.643401,
		0.656913,
		0.67071,
		0.684796,
		0.699179,
		0.713863,
		0.728855,
		0.744163,
		0.759792,
		0.775749,
		0.792042,
		0.808676,
		0.82566,
		0.843001,
		0.860706,
		0.878782,
		0.897239,
		0.916083,
		0.935322,
		0.954966,
		0.975022,
		0.9955,
		1.01641,
		1.03775,
		1.05955,
		1.0818,
		1.10452,
		1.12772,
		1.1514,
		1.17559,
		1.20028,
		1.22548,
		1.25122,
		1.2775,
		1.30433,
		1.33172,
		1.35969,
		1.38825,
		1.41741,
		1.44717,
		1.47757,
		1.5086,
		1.54028,
		1.57263,
		1.60566,
		1.63939,
		1.67382,
		1.70897,
		1.74486,
		1.78151,
		1.81892,
		1.85712,
		1.89613,
		1.93595,
		1.97661,
		2.01812,
		2.06051,
		2.10378,
		2.14797,
		2.19308,
		2.23914,
		2.28616,
		2.33418,
		2.3832,
		2.43325,
		2.48436,
		2.53653,
		2.58981,
		2.6442,
		2.69973,
		2.75643,
		2.81432,
		2.87343,
		2.93378,
		2.99539,
		3.0583,
		3.12253,
		3.18811,
		3.25507,
		3.32344,
		3.39323,
		3.4645,
		3.53726,
		3.61155,
		3.6874,
		3.76485,
		3.84392,
		3.92465,
		4.00707,
		4.09123,
		4.17715,
		4.26488,
		4.35445,
		4.44591,
		4.53928,
		4.63462,
		4.73195,
		4.83133,
		4.9328,
		5.0364,
		5.14218,
		5.25017,
		5.36044,
		5.47302,
		5.58796,
		5.70532,
		5.82515,
		5.94749,
		6.0724,
		6.19993,
		6.33014,
		6.46309,
		6.59883,
		6.73742,
		6.87892,
		7.02339,
		7.1709,
		7.3215,
		7.47527,
		7.63227,
		7.79256,
		7.95622,
		8.12332,
		8.29392,
		8.46812,
		8.64596,
		8.82755,
		9.01295,
		9.20224,
		9.3955,
		9.59283,
		9.7943,
		10.0,
	}
}

func CalcSI(resp []*Response) float64 {
	var si float64 = 0.0
	for i := 1; i < len(resp); i++ {
		r0 := resp[i-1]
		r1 := resp[i]
		if 0.1 < r1.Period && r1.Period <= 2.5 {
			si = si + (r0.Sv+r1.Sv)*(r1.Period-r0.Period)/2.0
		}
	}

	return si / 2.4
}
