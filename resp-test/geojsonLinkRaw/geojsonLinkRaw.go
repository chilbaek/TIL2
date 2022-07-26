package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type GeojsonFormat struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Crs  struct {
		Type       string `json:"type"`
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
	} `json:"crs"`
	Features []string `json:"features"`
}

func main() {
	lineCursor := 1

	file, err := os.Open("moct_node_sample.geojson")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	lineCount, _ := lineCounter(bufio.NewReader(file))

	file2, err := os.Open("moct_node_sample.geojson")
	fileScanner := bufio.NewScanner(file2)
	for fileScanner.Scan() {
		if (6 <= lineCursor) && (lineCursor < lineCount) {
			fmt.Println(fileScanner.Text())
		}
		lineCursor += 1
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	fmt.Printf("lineCursor: %d \n", lineCursor)

	file.Close()
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
