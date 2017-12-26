package response

import (
	"math"

	"github.com/takatoh/respspec/wave"
)

type Response struct {
	Freq float64
	Sa   float64
	Sv   float64
	Sd   float64
}

func NewResponse(freq, sa, sv, sd float64) *Response {
	p := new(Response)
	p.Freq = freq
	p.Sa = sa
	p.Sv = sv
	p.Sd = sd
	return p
}

func Resp(wave *wave.Wave, freq []float64, h float64) []*Response {
	var theta, dt, tdt, omega float64
	var ath, acd, abz float64
	var acc, vel, dis float64
	var a1, a2, a3 float64
	var am, f, vm, dm float64
	var k, c float64

	responses := make([]*Response, 0)
	theta = 1.4
	dt = wave.Dt / 10.0
	tdt = theta * dt
	nfreq := len(freq)
	z := interporate(wave.Data, 10)
	n := len(z)

	for j := 0; j < nfreq; j++ {
		if math.Abs(freq[j]) < 0.01 {
			am = 0.0
			for i := 1; i < n; i++ {
				if math.Abs(z[i]) > am {
					am = math.Abs(z[i])
				}
			}
			responses = append(responses, NewResponse(freq[j], am, 0.0, 0.0))
		} else {
			omega = 2.0 * math.Pi / freq[j]
			k = omega * omega
			c = 2.0 * h * omega

			// Constants for Willson's theta method.
			a1 = 1.0 + tdt * c / 2.0 + k * tdt * tdt / 6.0
			a2 = c + k * tdt
			a3 = tdt * c / 2.0 + k / 3.0 * tdt * tdt

			// Set initial values.
			acc = 0.0
			vel = 0.0
			dis = 0.0
			am = 0.0
			vm = 0.0
			dm = 0.0

			// Willson's theta method.
			for i := 1; i < n - 1; i++ {
				f = (theta - 1.0) * z[i] - theta * z[i + 1]
				ath = (f - k * dis - a2 * vel - a3 * acc) / a1
				acd = ((theta - 1.0) * acc + ath) / theta
				dis = dis + dt * vel + acc * dt * dt / 3.0 + acd * dt * dt / 6.0
				vel = vel + (acc + acd) * dt / 2.0
				acc = acd
				abz = math.Abs(acc + z[i + 1])
				if abz > am { am = abz }
				if math.Abs(vel) > vm { vm = math.Abs(vel) }
				if math.Abs(dis) > dm { dm = math.Abs(dis) }
			}

			responses = append(responses, NewResponse(freq[j], am, vm, dm))
		}
	}

	return responses
}

func interporate(zin []float64, ndiv int) []float64 {
	var k int
	var zinc float64
	nin := len(zin)
	z := make([]float64, 0)
	k = 0
	z = append(z, 0.0)
	for i := 0; i < nin - 1; i++ {
		if i == 0 {
			zinc = zin[i] / float64(ndiv)
		} else {
			zinc = (zin[i + 1] - zin[i]) / float64(ndiv)
		}
		for j := 0; j < ndiv; j++ {
			z = append(z, z[k] + zinc)
			k++
		}
	}

	return z
}
