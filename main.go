package main

import (
	"fmt"
	"os"

	"github.com/takatoh/respspec/wave"
)

func main() {
	csvfile := os.Args[1]

	wave := wave.LoadCSV(csvfile)

	fmt.Println(wave.Name)
	fmt.Println(wave.Dt)
	fmt.Println(len(wave.Data))
}
