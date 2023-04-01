package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/moeghifar/halaproxy/core"
)

func main() {
	fmt.Println("+++ halaproxy +++")
	config := core.Readflag()

	if config.DynamicBindNic.Enabled {
		transportsWithNic := []*http.Transport{}
		nics, err := net.Interfaces()
		if err != nil {
			panic(err)
		}

		// get custom nic
		for _, nic := range nics {
			if strings.HasPrefix(nic.Name, config.DynamicBindNic.Prefix) {
				transportsWithNic = append(
					transportsWithNic,
					core.CreateTransportWithNic(nic),
				)
			}
		}

		startingPort := config.DynamicBindNic.StartPort

		for _, nicTransport := range transportsWithNic {
			go listenAndServe(startingPort, core.SetServer(nicTransport))
			startingPort++
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
