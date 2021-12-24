package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"unsafe"
)

var (
	Port = "2125"
)

func fileTravelerServer() {
	l, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		fmt.Println("Failed to listen port", Port, "error:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("File traveler server is started on port", Port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection, error:", err.Error())
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	bufHeader := make([]byte, unsafe.Sizeof(FileHeader{}))
	_, err := conn.Read(bufHeader)
	if err != nil {
		fmt.Println("Read file header error, error:", err.Error())
		return
	}
	header := (*FileHeader)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&bufHeader)).Data))

	bufFileName := make([]byte, header.FileNameLength)
	_, err = conn.Read(bufFileName)
	if err != nil {
		fmt.Println("Read file name error, error:", err.Error())
		return
	}
	fileName := string(bufFileName)
	fmt.Println("Receiving new file, file name:", fileName)

	targetFile, err := os.Create("./" + fileName)
	if err != nil {
		fmt.Println("Create target file error, error:", err.Error())
		return
	}

	progChan := make(chan int)
	maxProg := ProgressBarLength

	go func() {
		buf := make([]byte, BufferSize)
		currBytes := 0
		lastProgress := 0
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("Read file error, error:", err.Error())
				return
			}
			currBytes += n
			newProgress := int(float64(currBytes) / float64(header.FileLength) * float64(maxProg))
			if newProgress-lastProgress >= 1 {
				progChan <- newProgress
			}
			if n == 0 {
				close(progChan)
				break
			}
			_, err = targetFile.Write(buf[:n])
			if err != nil && err != io.EOF {
				fmt.Println("Write file error, error:", err.Error())
				return
			}
		}
	}()

	progressBar(maxProg, progChan)

	fmt.Println("Finished receiving file", fileName)
}
