package wave

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
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

	read_file, _ := os.OpenFile(filename, os.O_RDONLY, 0600)
	reader = csv.NewReader(read_file)

	columns, err = reader.Read()
	wave.Name = columns[1]
	for {
		columns, err = reader.Read()
		if err == io.EOF {
			wave.Dt = t2 - t1
			wave.Data = data
			return wave
		}
		t1 = t2
		t2, _ = strconv.PaarseFloat(columns[0], 64)
		d = strconv.PaarseFloat(columns[1], 64)
		data = append(data, d)
	}
}
