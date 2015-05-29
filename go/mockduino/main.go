/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Christian Speckner
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

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
