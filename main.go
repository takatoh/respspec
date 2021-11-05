package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/takatoh/respspec/response"
	"github.com/takatoh/seismicwave"
)

const (
	progVersion = "v0.8.0"
)

func main() {
	var waves []*seismicwave.Wave
	var err error
	var period []float64

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage:
  %s [options] <wavefile>

Options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	opt_period := flag.String("period", "", "Specify period file.")
	opt_h := flag.Float64("h", 0.05, "Specify attenuation constant.")
	opt_max := flag.Float64("max", 0.0, "Specify maximum acc.")
	opt_format := flag.String("format", "", "wave format.")
	opt_name := flag.String("name", "unnamed", "wave name.")
	opt_dt := flag.Float64("dt", 0.0, "time delta.")
	opt_ndata := flag.Int("ndata", 0, "number of data.")
	opt_skip := flag.Int("skip", 0, "skip lines.")
	opt_si := flag.Bool("si", false, "Calculate SI.")
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
	if *opt_si {
		vals := []float64{0.1, 2.5}
		period = insertPeriod(period, vals)
	}

	srcfile := flag.Args()[0]

	if *opt_format != "" {
		if *opt_dt == 0.0 || *opt_ndata == 0 {
			fmt.Fprintln(os.Stderr, "Error: At least -dt and -ndata option must be given.")
			os.Exit(1)
		}
		waves, err = seismicwave.LoadFixedFormat(srcfile, *opt_name, *opt_format, *opt_dt, *opt_ndata, *opt_skip)
	} else {
		waves, err = seismicwave.LoadCSV(srcfile)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	wv := waves[0]

	if *opt_max > 0.0 {
		max := wv.AbsMax()
		wv = mul(wv, *opt_max/max)
	}

	spectrum := response.Spectrum(wv, period, *opt_h)

	fmt.Println(wv.Name)
	if *opt_si {
		si := response.CalcSI(spectrum)
		fmt.Printf("SI = %f\n", si)
	} else {
		fmt.Println("Period,Sa,Sv,Sd")
		for _, res := range spectrum {
			fmt.Printf("%f,%f,%f,%f\n", res.Period, res.Sa, res.Sv, res.Sd)
		}
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

func mul(w *seismicwave.Wave, fac float64) *seismicwave.Wave {
	data := w.Data
	for i := 0; i < len(data); i++ {
		data[i] *= fac
	}
	w.Data = data

	return w
}

func insertPeriod(period []float64, vals []float64) []float64 {
	for _, x := range vals {
		if !(find(period, x)) {
			period = append(period, x)
		}
	}
	sort.Float64s(period)

	return period
}

func find(s []float64, x float64) bool {
	for _, v := range s {
		if x == v {
			return true
		}
	}
	return false
}
