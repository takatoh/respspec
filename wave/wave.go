package wave

import (
	"encoding/csv"
	"io"
	"bufio"
	"os"
	"strings"
	"strconv"
	"math"
)

type Wave struct {
	Name string
	Dt   float64
	Data []float64
}

func NewWave() *Wave {
	p := new(Wave)
	return p
}

func (w *Wave) AbsMax() float64 {
	max := 0.0
	for _, d := range w.Data {
		if math.Abs(d) > max {
			max = math.Abs(d)
		}
	}
	return max
}

func (w *Wave) Mul(factor float64) *Wave {
	n := len(w.Data)
	data := make([]float64, n)
	for i := 0; i < n; i++ {
		data[i] = w.Data[i] * factor
	}
	w.Data = data
	return w
}

func LoadCSV(filename string) *Wave {
	var reader *csv.Reader
	var columns []string
	var err error
	var wave *Wave
	var t1, t2, d float64
	var data []float64

	wave = NewWave()
	t1 = 0.0
	t2 = 0.0

	read_file, _ := os.Open(filename)
	reader = csv.NewReader(read_file)

	columns, err = reader.Read()
	wave.Name = columns[1]
	for {
		columns, err = reader.Read()
		if err == io.EOF {
			wave.Dt = round(t2 - t1, 2)
			wave.Data = data
			return wave
		}
		t1 = t2
		t2, _ = strconv.ParseFloat(columns[0], 64)
		d, _ = strconv.ParseFloat(columns[1], 64)
		data = append(data, d)
	}
}

func round(val float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}

func LoadWave(filename, name, format string, dt float64, n, skip int) *Wave {
	wave := NewWave()
	wave.Name = name
	wave.Dt = dt

	d_num, d_len := parseFormat(format)
	line_num := int(math.Ceil(float64(n) / float64(d_num)))
	data := make([]string, 0)

	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for i := 0; i < skip; i++ {
		scanner.Scan()
	}
	for i := 0; i < line_num; i++ {
		scanner.Scan()
		line := scanner.Text()
		runes := []rune(line)
		for j := 0; j < d_num; j++ {
			data = append(data, string(runes[j * d_len:(j + 1) * d_len]))
		}
	}

	for i := 0; i < n; i++ {
		d, _ := strconv.ParseFloat(strings.Trim(data[i], " "), 64)
		wave.Data = append(wave.Data, d)
	}

	return wave
}

func parseFormat(format string) (int, int) {
	ss := strings.Split(format, "F")
	i64, _ := strconv.ParseInt(ss[0], 10, 64)
	d_num := int(i64)
	ss = strings.Split(ss[1], ".")
	i64, _ = strconv.ParseInt(ss[0], 10, 64)
	d_len := int(i64)
	return d_num, d_len
}
