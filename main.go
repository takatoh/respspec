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
	progVersion = "v0.4.0"
)

func main() {
	var wv *wave.Wave
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
	opt_max := flag.Float64("max", 0.0, "Specify maximum acc.")
	opt_format := flag.String("format", "", "wave format.")
	opt_name := flag.String("name", "unnamed", "wave name.")
	opt_dt := flag.Float64("dt", 0.0, "time delta.")
	opt_num := flag.Int("num", 0, "number of data.")
	opt_skip := flag.Int("skip", 0, "skip lines.")
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

	srcfile := flag.Args()[0]
	if *opt_format != "" {
		if *opt_dt == 0.0 || *opt_num == 0 {
			fmt.Fprintln(os.Stderr, "Error: At least -dt and -num option must be given.")
			os.Exit(1)
		}
		wv = wave.LoadWave(srcfile, *opt_name, *opt_format, *opt_dt, *opt_num, *opt_skip)
	} else {
		wv = wave.LoadCSV(srcfile)
	}

	if *opt_max > 0.0 {
		max := wv.AbsMax()
		wv = wv.Mul(*opt_max / max)
	}

	responses := response.Resp(wv, freq, h)

	fmt.Println(wv.Name)
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
