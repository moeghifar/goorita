package core

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type config struct {
	HttpPort       int
	DynamicBindNic DynamicBindNic
}

type DynamicBindNic struct {
	Enabled   bool
	Prefix    string
	StartPort int
}

var defaultDynamicBindNicStartPort = 33001
var defaultHttpPort = 33000

func Readflag() config {
	var httpPort int
	var dynamicBindNic string

	flag.IntVar(&httpPort, "http-port", 33000, "cli command")
	flag.StringVar(&dynamicBindNic, "dynamic-nic-bind", "", "cli command")
	flag.Parse()

	cfg := config{
		HttpPort: defaultHttpPort,
		DynamicBindNic: DynamicBindNic{
			Enabled:   false,
			StartPort: defaultDynamicBindNicStartPort,
			Prefix:    "*",
		},
	}

	if httpPort > 0 {
		cfg.HttpPort = httpPort
	}

	if dynamicBindNic != "" {
		cfg.DynamicBindNic = parseDynamicBindNic(dynamicBindNic)
	}

	return cfg
}

func parseDynamicBindNic(input string) DynamicBindNic {
	split := strings.Split(input, ":")

	if len(split) < 2 {
		panic("invalid format --dynamic-bind-port example [enx:330001]")
	}

	startPort := defaultDynamicBindNicStartPort

	valid, err := strconv.Atoi(split[1])
	if err != nil {
		panic(fmt.Sprintf("%+v: invalid format --dynamic-bind-port example [enx:330001]", err))
	}

	if valid > 0 {
		startPort = valid
	}

	return DynamicBindNic{
		Enabled:   true,
		Prefix:    split[0],
		StartPort: startPort,
	}
}
