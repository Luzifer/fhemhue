package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	httpHelper "github.com/Luzifer/go_helpers/http"
	"github.com/Luzifer/rconfig"
	log "github.com/Sirupsen/logrus"
	"github.com/pborges/huemulator"
	"github.com/satori/uuid"
)

var (
	cfg struct {
		ConfigFile     string `flag:"config,c" default:"config.yml" description:"Configuration file with FHEM details"`
		Listen         string `flag:"listen" default:"127.0.0.1:10000" description:"IP/Port to listen on (IP is required and cannot be 0.0.0.0)"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Print version information and exit"`
	}

	fileConfig *config

	version = "dev"
)

func init() {
	var err error
	if err = rconfig.Parse(&cfg); err != nil {
		log.Fatalf("Unable to parse CLI arguments: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("fhemhue %s\n", version)
		os.Exit(0)
	}

	fileConfig, err = loadConfig(cfg.ConfigFile)
	if err != nil {
		log.Fatalf("Unable to load configuration file: %s", err)
	}
}

func main() {
	ip, port, err := net.SplitHostPort(cfg.Listen)
	if err != nil {
		log.Fatalf("Looks like your listen address %q is invalid: %s", cfg.Listen, err)
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Port %q does not look like a port number: %s", port, err)
	}

	hueConfig := huemulator.Config{
		Hostname: ip,
		Port:     portNum,
		UDN:      uuid.NewV5(uuid.NamespaceOID, "fhemhue").String(),
		Protocol: "http",
	}

	go huemulator.UpnpResponder(hueConfig)
	r, _ := huemulator.NewRouter(hueConfig, fileConfig)
	log.Fatalf("HTTP Server stopped: %s", http.ListenAndServe(cfg.Listen, httpHelper.NewHTTPLogHandler(r)))
}
