package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/moeghifar/halaproxy/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.NewConsoleWriter()

	log.Info().Msg("+++ halaproxy +++")

	config := core.Readflag()

	if config.DynamicBindNic.Enabled {
		startingPort := config.DynamicBindNic.StartPort
		nics, err := net.Interfaces()
		if err != nil {
			panic(err)
		}

		for _, nic := range nics {
			if strings.HasPrefix(nic.Name, config.DynamicBindNic.Prefix) {
				fmt.Println("found nic", nic.Name, "bind to port", startingPort)
				go listenAndServe(startingPort, core.SetServer(core.CreateTransportWithNic(nic)))
				startingPort++
			}
		}
	}

	listenAndServe(config.HttpPort, core.SetServer(core.CreateTransport(nil)))
}

func listenAndServe(port int, srv http.Handler) {
	fmt.Println("serving in port", port)
	http.ListenAndServe(intToPort(port), srv)
}

func intToPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
