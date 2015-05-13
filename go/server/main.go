package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/server/controls"
	"github.com/DirtyHairy/arduino-butler/go/util/ip"
	routerPkg "github.com/DirtyHairy/arduino-butler/go/util/router"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type configBag struct {
	port         int
	ip           string
	frontendPath string
	configFile   string
}

func parseCommandline() configBag {
	config := configBag{
		port:         8888,
		frontendPath: "./frontend",
		configFile:   "config.json",
	}

	ip := ip.Create()

	flag.Var(&ip, "i", "server listen IP")
	flag.IntVar(&config.port, "p", config.port, "server listen port")
	flag.StringVar(&config.frontendPath, "f", config.frontendPath, "path to frontend")
	flag.StringVar(&config.configFile, "c", config.configFile, "configuration file")

	flag.Parse()

	config.ip = ip.String()

	return config
}

func createContrlSetFromConfigFile(configFile string) (*controls.ControlSet, error) {
	configFileContent, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	var controlSetMarshalled controls.MarshalledControlSet

	err = json.Unmarshal(configFileContent, &controlSetMarshalled)

	if err != nil {
		return nil, err
	}

	controlSet, err := controlSetMarshalled.Unmarshal()

	if err != nil {
		return nil, err
	}

	return controlSet, nil
}

func main() {
	config := parseCommandline()

	listenAddress := config.ip + ":" + strconv.Itoa(config.port)

	fmt.Println("reading config...")

	controlSet, err := createContrlSetFromConfigFile(config.configFile)

	if err != nil {
		fmt.Println(err)
		return
	}

	controller := CreateController(controlSet)

	router := routerPkg.CreateRouter(10)
	router.AddRoute("^/api/switch/(\\w+)/(on|off)$", routerPkg.HandlerFunction(controller.HandleSwitch))
	router.AddRoute("^/api/structure$", routerPkg.HandlerFunction(controller.GetStructure))

	http.Handle("/", http.FileServer(http.Dir(config.frontendPath)))
	http.Handle("/api/", router)

	fmt.Printf("Frontend served from %s\n", config.frontendPath)
	fmt.Printf("Server listening on %s\n", listenAddress)

	if err := controlSet.Start(); err != nil {
		panic(fmt.Sprintf("failed to start controls: %v", err))
	}
	defer controlSet.Stop()

	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
