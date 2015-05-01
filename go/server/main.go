package main

import (
	"flag"
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/util/ip"
	"net/http"
	"os"
	"strconv"
)

type configBag struct {
	port         int
	ip           string
	frontendPath string
}

func parseCommandline() configBag {
	config := configBag{8888, "", "./frontend"}

	ip := ip.Create()

	flag.Var(&ip, "i", "server listen IP")
	flag.IntVar(&config.port, "p", 8888, "server listen port")
	flag.StringVar(&config.frontendPath, "f", "./frontend", "path to frontend")
	flag.Parse()

	config.ip = ip.String()

	return config
}

func main() {
	config := parseCommandline()

	listenAddress := config.ip + ":" + strconv.Itoa(config.port)

	http.Handle("/", http.FileServer(http.Dir(config.frontendPath)))

	fmt.Printf("Frontend served from %s\n", config.frontendPath)
	fmt.Printf("Server listening on %s\n", listenAddress)

	controlChan := make(chan error)

	go func() {
		if err := http.ListenAndServe(listenAddress, nil); err != nil {
			controlChan <- err
		} else {
			controlChan <- nil
		}
	}()

	if err := <-controlChan; err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
