package main

import (
	"flag"
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/util"
	"github.com/DirtyHairy/arduino-butler/go/util/ip"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type configBag struct {
	port  int
	ip    string
	delay time.Duration
}

var switchRouteRx *regexp.Regexp
var delay time.Duration
var mutex sync.Mutex

func parseCommandline() configBag {
	config := configBag{
		port:  8000,
		ip:    "",
		delay: 200 * time.Millisecond,
	}

	ip := ip.Create()
	delay := util.DurationValue(config.delay)

	flag.Var(&ip, "i", "server listen IP")
	flag.Var(&delay, "d", "request delay")

	flag.IntVar(&config.port, "p", 8000, "server listen port")
	flag.Parse()

	config.ip = ip.String()
	config.delay = delay.Value()

	return config
}

func setupHeaders(response http.ResponseWriter) {
	headers := response.Header()

	headers.Add("cache-control", "no-cache")
	headers.Add("access-control-allow-origin", "*")
}

func send404(response http.ResponseWriter) {
	response.WriteHeader(http.StatusNotFound)
	response.Write([]byte("<html><head><title>Page not found!</title></head><body>Page not found!</body></html>\n"))
}

func handleSwitch(response http.ResponseWriter, request *http.Request) {
	mutex.Lock()

	defer func() {
		mutex.Unlock()
	}()

	setupHeaders(response)

	matches := switchRouteRx.FindStringSubmatch(request.URL.Path)

	if matches == nil {
		send404(response)
		return
	}

	socketIdx, _ := strconv.Atoi(matches[1])

	time.Sleep(delay)

	fmt.Printf("Toggled socket %d %s\n", socketIdx, matches[2])

	response.WriteHeader(http.StatusOK)
}

func init() {
	switchRouteRx = util.CompileRegex("^/socket/(\\d+)/(on|off)$")
}

func main() {
	config := parseCommandline()
	listenAddress := config.ip + ":" + strconv.Itoa(config.port)
	delay = config.delay

	fmt.Printf("response delay: %v\n", delay)
	fmt.Printf("Server listening on %s...\n", listenAddress)
	http.HandleFunc("/", handleSwitch)

	err := http.ListenAndServe(listenAddress, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
