package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"unsafe"
)

func fileTravelerClient(filePath string, targetHostName string) {
	conn, err := net.Dial("tcp", targetHostName+":"+Port)
	if err != nil {
		fmt.Println("Failed to connect to host", targetHostName, "error:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file", filePath, "error:", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	fileStat, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("Failed to stat file", filePath, "error:", err.Error())
		os.Exit(1)
	}

	fileName := fileStat.Name()
	fileNameLength := uint64(len(fileName))
	fileLength := fileStat.Size()

	_, err = conn.Write((*(*[8]byte)(unsafe.Pointer(&fileNameLength)))[:])
	if err != nil {
		fmt.Println("Failed to write file name length", filePath, "error:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write((*(*[8]byte)(unsafe.Pointer(&fileLength)))[:])
	if err != nil {
		fmt.Println("Failed to write file length", filePath, "error:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Failed to write file name", filePath, "error:", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Failed to read file", filePath, "error:", err.Error())
			os.Exit(1)
		}
		if n == 0 {
			break
		}
		_, err = conn.Write(buf[:n])
		if err != nil && err != io.EOF {
			fmt.Println("Failed to write data to target", targetHostName, "error:", err.Error())
			os.Exit(1)
		}
	}

	fmt.Println("File", filePath, "is sent successfully!")
}
