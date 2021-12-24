package main

import "fmt"

const (
	BufferSize        = 8 * 1024
	ProgressBarLength = 25
)

type FileHeader struct {
	FileNameLength uint64
	FileLength     uint64
}

func progressBar(max int, progress chan int) {
	fmt.Print("[")
	for i := 0; i < max; i++ {
		fmt.Print("-")
	}
	fmt.Print("] - 0%")
	for prog := range progress {
		fmt.Print("\r[")
		for i := 0; i < prog; i++ {
			fmt.Print("=")
		}
		for i := 0; i < max-prog; i++ {
			fmt.Print("-")
		}
		fmt.Printf("] - %d%%", int64(float64(prog)/float64(max)*100))
		if prog >= max {
			break
		}
	}
	fmt.Println()
}
