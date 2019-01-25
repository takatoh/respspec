package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"flag"

	"github.com/takatoh/seismicwave"
	"github.com/takatoh/respspec/response"
)

const (
	progVersion = "v0.5.0"
)

func main() {
	var waves []*seismicwave.Wave
	var err error
	var period []float64
	var h float64 = 0.05

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
`Usage:
  %s [options] <wavefile>

Options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	opt_period := flag.String("period", "", "Specify period file.")
	opt_max := flag.Float64("max", 0.0, "Specify maximum acc.")
	opt_format := flag.String("format", "", "wave format.")
	opt_name := flag.String("name", "unnamed", "wave name.")
	opt_dt := flag.Float64("dt", 0.0, "time delta.")
	opt_ndata := flag.Int("ndata", 0, "number of data.")
	opt_skip := flag.Int("skip", 0, "skip lines.")
	opt_version := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}
	if *opt_period != "" {
		period = loadPeriod(*opt_period)
	} else {
		period = response.DefaultPeriod()
	}

	srcfile := flag.Args()[0]
//	fp, err := os.Open(srcfile)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Cannot open file: %s\n", srcfile)
//		os.Exit(1)
//	}
//	defer fp.Close()

	if *opt_format != "" {
		if *opt_dt == 0.0 || *opt_ndata == 0 {
			fmt.Fprintln(os.Stderr, "Error: At least -dt and -ndata option must be given.")
			os.Exit(1)
		}
		waves, err = seismicwave.LoadFixedFormat(srcfile, *opt_name, *opt_format, *opt_dt, *opt_ndata, *opt_skip)
	} else {
		waves, err = seismicwave.LoadCSV(srcfile)
	}
	wv := waves[0]

//	if *opt_max > 0.0 {
//		max := wv.AbsMax()
//		wv = wv.Mul(*opt_max / max)
//	}

	responses := response.Resp(wv, period, h)

	fmt.Println(wv.Name)
	fmt.Println("Period,Sa,Sv,Sd")
	for _, res := range responses {
		fmt.Printf("%f,%f,%f,%f\n", res.Period, res.Sa, res.Sv, res.Sd)
	}
}

func loadPeriod(filename string) []float64 {
	period := make([]float64, 0)

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
		period = append(period, v)
	}

	return period
}
