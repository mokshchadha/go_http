package main

import (
	"fmt"
	"io"
	"log"
	"net"
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

	listener, err := net.Listen("tcp4", ":42069")

	if err != nil {
		log.Println("Could not register listener to TCP", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("could not connect with the server", err)
			continue
		}

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}

}
