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
			log.Panic().Err(err)
		}

		for _, nic := range nics {
			if strings.HasPrefix(nic.Name, config.DynamicBindNic.Prefix) {
				log.Info().Msgf("found nic %s bind to port %d", nic.Name, startingPort)
				go func() {
					if err := listenAndServe(startingPort, core.SetServer(core.CreateTransportWithNic(nic))); err != nil {
						log.Panic().Err(err)
					}
				}()
				startingPort++
			}
		}
	}

	if err := listenAndServe(config.HttpPort, core.SetServer(core.CreateTransport(nil))); err != nil {
		log.Fatal().Err(err)
	}
}

func listenAndServe(port int, handler http.Handler) error {
	log.Info().Msgf("serving in port %d", port)
	return http.ListenAndServe(intToPort(port), handler)
}

func intToPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
