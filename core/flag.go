package core

import (
	"flag"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
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
		log.Fatal().Msg("invalid format --dynamic-bind-port example [enx:330001]")
	}

	startPort := defaultDynamicBindNicStartPort

	valid, err := strconv.Atoi(split[1])
	if err != nil {
		log.Fatal().Msgf("%+v: invalid format --dynamic-bind-port example [enx:330001]", err)
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
