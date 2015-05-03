package main

import (
	"flag"
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/server/controls"
	"github.com/DirtyHairy/arduino-butler/go/util/ip"
	routerPkg "github.com/DirtyHairy/arduino-butler/go/util/router"
	"net/http"
	"os"
	"strconv"
)

type configBag struct {
	port         int
	ip           string
	frontendPath string
	controlHost  string
}

func parseCommandline() configBag {
	config := configBag{8888, "", "./frontend", "localhost:8080"}

	ip := ip.Create()

	flag.Var(&ip, "i", "server listen IP")
	flag.IntVar(&config.port, "p", config.port, "server listen port")
	flag.StringVar(&config.frontendPath, "f", config.frontendPath, "path to frontend")
	flag.StringVar(&config.controlHost, "h", config.controlHost, "control host")

	flag.Parse()

	config.ip = ip.String()

	return config
}

func handleSwitch(response http.ResponseWriter, request *http.Request, matches []string) {
	fmt.Printf("handling request for %s\n", matches[0])

	switchIndex, _ := strconv.Atoi(matches[1])
	toggle := matches[2] == "on"

	err := controls.SendSwitchCommand(controls.SwitchCommand{
		Index:  uint(switchIndex),
		Toggle: toggle,
	})

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	switch err.(type) {
	case nil:
		response.WriteHeader(http.StatusOK)

	case controls.SwitchNotFoundError:
		response.WriteHeader(http.StatusNotFound)

	default:
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	config := parseCommandline()

	listenAddress := config.ip + ":" + strconv.Itoa(config.port)

	router := routerPkg.CreateRouter(10)
	router.AddRoute("^/api/switch/(\\d)/(on|off)$", routerPkg.HandlerFunction(handleSwitch))

	http.Handle("/", http.FileServer(http.Dir(config.frontendPath)))
	http.Handle("/api/", router)

	fmt.Printf("Frontend served from %s\n", config.frontendPath)
	fmt.Printf("Server listening on %s\n", listenAddress)

	controls.StartCommandProcessor(config.controlHost)
	defer controls.StopCommandProcessor()

	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
