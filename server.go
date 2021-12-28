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
	bufHeader := make([]byte, int(unsafe.Sizeof(FileHeader{})))
	_, err := io.ReadFull(conn, bufHeader)
	if err != nil {
		fmt.Println("Read file header error, error:", err.Error())
		return
	}
	header := (*FileHeader)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&bufHeader)).Data))

	bufFileName := make([]byte, int(header.FileNameLength))
	_, err = io.ReadFull(conn, bufFileName)
	if err != nil {
		fmt.Println("Read file name error, error:", err.Error())
		return
	}
	fileName := string(bufFileName)

	targetFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Create target file error, error:", err.Error())
		return
	}
	defer targetFile.Close()

	progChan := make(chan int)
	syncChan := make(chan bool)
	maxProg := ProgressBarLength

	go func() {
		progressBar(maxProg, progChan)
		fmt.Println("Finished receiving file", fileName)
		syncChan <- true
	}()

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
		lastProgress = newProgress
		if err == io.EOF {
			break
		}
		_, err = targetFile.Write(buf[:n])
		if err != nil && err != io.EOF {
			fmt.Println("Write file error, error:", err.Error())
			return
		}
	}
	<-syncChan
}
