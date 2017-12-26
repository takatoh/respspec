package response

import (
	"math"
)

type Response struct {
	Freq float64
	Sa   float64
	Sv   float64
	Sd   float64
}

func NewResponse(freq float64) *Response {
	p := new(Response)
	p.Freq = freq
	return p
}


