package main

import (
	"flag"

	"github.com/flashbots/boost-relay/apis/registervalidator"
	"github.com/sirupsen/logrus"
)

var (
	version = "dev" // is set during build process

	// defaults
	defaultListenAddr     = "localhost:9064"
	defaultBeaconEndpoint = "http://localhost:5052"

	// cli flags
	listenAddr     = flag.String("listen-addr", defaultListenAddr, "listen address")
	beaconEndpoint = flag.String("beacon-endpoint", defaultBeaconEndpoint, "beacon endpoint")
)

var log = logrus.WithField("module", "cmd/relay")

func main() {
	flag.Parse()
	log.Printf("boost-relay %s [registervalidator-api]", version)

	validatorService := registervalidator.NewBeaconClientValidatorService(*beaconEndpoint)
	// TODO: should be done at the start of every epoch
	err := validatorService.FetchValidators()
	if err != nil {
		log.WithError(err).Fatal("failed to fetch validators from beacon node")
	}

	srv, err := registervalidator.NewRegisterValidatorService(*listenAddr, validatorService, log)
	if err != nil {
		log.WithError(err).Fatal("failed to create service")
	}

	log.Println("listening on", *listenAddr)
	log.Fatal(srv.StartHTTPServer())
}
