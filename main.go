package main

import (
	"fmt"
	"os"
//	"io"
	"bufio"
	"strconv"
	"flag"

	"github.com/takatoh/respspec/wave"
	"github.com/takatoh/respspec/response"
)

const (
	progVersion = "v0.1.0"
)

func main() {
	var freq []float64
	var defaultFreq []float64 = []float64{
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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
`Usage:
  %s [options] <file.csv>
options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	opt_freq := flag.String("freq", "", "Specify calcurate frequency.")
	opt_version := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}
	if *opt_freq != "" {
		freq = loadFreq(*opt_freq)
	} else {
		freq = defaultFreq
	}

	csvfile := flag.Args()[0]

	wave := wave.LoadCSV(csvfile)
	responses := response.Resp(wave, freq, h)

	fmt.Println(wave.Name)
	fmt.Println("Freq,Sa,Sv,Sd")
	for _, res := range responses {
		fmt.Printf("%f,%f,%f,%f\n", res.Freq, res.Sa, res.Sv, res.Sd)
	}
}

func loadFreq(filename string) []float64 {
	freq := make([]float64, 0)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
//		fmt.Println(s)
		v, _ := strconv.ParseFloat(s, 64)
		freq = append(freq, v)
	}

	return freq
}
