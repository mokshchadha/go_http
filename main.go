package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	myChan := make(chan string)

	go func() {
		str := ""

		defer close(myChan)
		defer f.Close()
		for {
			data := make([]byte, 8)

			_, err := f.Read(data)

			if err != nil {
				if err == io.EOF {
					myChan <- str
					break
				}
			}
			str += string(data)
			if strings.Contains(string(str), "\n") {
				splits := strings.Split(string(str), "\n")
				line, remaining := splits[0], splits[1]
				str = remaining
				myChan <- line
			}

		}

	}()
	return myChan
}

func main() {
	file, err := os.Open("messages.txt")

	if err != nil {
		log.Fatal("Could not read the file", err)
	}

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
