package main

import (
	"fmt"
	"io"
	"os"
)

// this is a comment

func readFile(filename string) string {
	//fmt.Println(filename)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	text := ""
	for {
		buf := make([]byte, 1024)
		// read a chunk
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		text += fmt.Sprintf("%s", buf)
	}
	return text
}

func main() {
	if len(os.Args) <= 1 {
		os.Exit(2)
	}
	html := readFile(os.Args[1])
}