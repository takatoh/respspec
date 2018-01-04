package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"flag"

	"github.com/takatoh/respspec/wave"
	"github.com/takatoh/respspec/response"
)

const (
	progVersion = "v0.2.1"
)

func main() {
	var freq []float64
	var h float64 = 0.05

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
`Usage:
  %s [options] <file.csv>

Options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	opt_freq := flag.String("freq", "", "Specify frequency file.")
	opt_version := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}
	if *opt_freq != "" {
		freq = loadFreq(*opt_freq)
	} else {
		freq = response.DefaultFreq()
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
		v, _ := strconv.ParseFloat(s, 64)
		freq = append(freq, v)
	}

	return freq
}
