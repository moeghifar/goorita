package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/moeghifar/meng/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.NewConsoleWriter()

	log.Info().Msg("+++ meng +++")

	config := core.Readflag()

	if config.DynamicBindNic.Enabled {
		startingPort := config.DynamicBindNic.StartPort
		nics, err := net.Interfaces()
		if err != nil {
			panic(err)
		}

		for _, nic := range nics {
			if strings.HasPrefix(nic.Name, config.DynamicBindNic.Prefix) {
				log.Info().Msgf("found nic %s bind to port %d", nic.Name, startingPort)
				go listenAndServe(startingPort, core.SetServer(core.CreateTransportWithNic(nic)))
				startingPort++
			}
		}
	}

	if err := listenAndServe(config.HttpPort, core.SetServer(core.CreateTransport(nil))); err != nil {
		panic(err)
	}
}

func listenAndServe(port int, srv http.Handler) error {
	log.Info().Msgf("serving in port %d", port)
	return http.ListenAndServe(intToPort(port), srv)
}

func intToPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
