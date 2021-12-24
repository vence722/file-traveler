package main

import (
	"fmt"
	"os"
	"strings"
)

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
		fmt.Printf("\r["+strings.Repeat("=", prog)+strings.Repeat("-", max-prog)+"] - %d%%", int64(float64(prog)/float64(max)*100))
		os.Stdout.Sync()
		if prog >= max {
			// ensure the last progress bar shows correct progress info
			fmt.Println("\r[" + strings.Repeat("=", max) + "] - 100%")
			break
		}
	}
}
