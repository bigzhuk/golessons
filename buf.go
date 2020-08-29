package main

import (
	"bufio"
	"fmt"
	"os"
)

type myFile struct {
	counts map[string]int
}

func main() {
	myFiles := make(map[string]myFile)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, myFiles, "Stdin")
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, myFiles, arg)
			f.Close()
		}
	}

	for name, myFile := range myFiles {
		fmt.Println(name)
		for line, n := range myFile.counts {
			if n > 1 {
				fmt.Printf("%d\t%s\n", n, line)
			}
		}
	}
}
func countLines(f *os.File, myFiles map[string]myFile, filename string) {
	input := bufio.NewScanner(f)
	file := new(myFile)
	file.counts = make(map[string]int)
	for input.Scan() {
		if input.Text() == "" {
			break
		}
		file.counts[input.Text()]++
		myFiles[filename] = *file
	}
}
