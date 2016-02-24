package main

import (
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"time"
)

func readLine(reader io.Reader) (line string, err error) {
	readBuffer := make([]byte, 1, 1)

	for {
		var n int

		n, err = reader.Read(readBuffer)
		if err != nil && err != io.EOF {
			return
		}

		if n != 1 {
			err = errors.New("read timeout --- handshake already complete?")
			return
		}

		if readBuffer[0] == 0 || readBuffer[0] == 0x0A || readBuffer[0] == 0x0D {
			return
		}

		if n == 1 {
			line = line + string(readBuffer)
		}
	}
}

func transmitLine(buffer []byte, out io.Writer) (err error) {
	var transmitCount int

	for transmitCount != len(buffer) {
		var partialCount int

		partialCount, err = out.Write(buffer[transmitCount:])
		if err != nil {
			return
		}

		transmitCount = transmitCount + partialCount
	}

	transmitCount = 0
	for transmitCount == 0 && err == nil {
		transmitCount, err = out.Write([]byte{0x0A})
	}

	return
}

func waitFor(key string, reader io.Reader) (err error) {
	var line string

	for line != key {
		line, err = readLine(reader)

		if err != nil {
			return
		}
	}

	return
}

func waitForWithTimeout(key string, reader io.Reader, timeout time.Duration) (success bool, err error) {
	var line string
	now := time.Now()

	for line != key && time.Now().Sub(now) < timeout {
		line, err = readLine(reader)

		if err != nil {
			return
		}
	}

	success = line == key
	return
}

func main() {
	config := serial.Config{
		Name:        "/dev/ttyUSB0",
		Baud:        115200,
		ReadTimeout: 5 * time.Second,
	}

	s, err := serial.OpenPort(&config)
	if err != nil {
		log.Fatal(err)
	}

	var handshakeComplete bool

	for !handshakeComplete {
		fmt.Println("WAITING")

		err = waitFor("waiting", s)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("CONNECT")

		err = transmitLine([]byte("connect"), s)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("READY")

		handshakeComplete, err = waitForWithTimeout("ready", s, 2*time.Second)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("handshake complete!")
}
