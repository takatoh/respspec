package main

import (
	"fmt"
	"os"
	"flag"

	"github.com/takatoh/respspec/wave"
	"github.com/takatoh/respspec/response"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
`Usage:
  %s [options] <file.csv>
options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	csvfile := flag.Args()[0]

	wave := wave.LoadCSV(csvfile)

	var freq []float64 = []float64{
		0.0,
		0.1,
		0.2,
		0.3,
		0.4,
		0.5,
		0.6,
		0.7,
		0.8,
		0.9,
		1.0,
		1.1,
		1.2,
		1.3,
		1.4,
		1.5,
		1.6,
		1.7,
		1.8,
		1.9,
		2.0,
		2.5,
		3.0,
		4.0,
		5.0,
	}
	var h float64 = 0.05

	responses := response.Resp(wave, freq, h)

	fmt.Println(wave.Name)
	fmt.Println("Freq,Sa,Sv,Sd")
	for _, res := range responses {
		fmt.Printf("%f,%f,%f,%f\n", res.Freq, res.Sa, res.Sv, res.Sd)
	}
}
