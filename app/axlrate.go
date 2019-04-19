package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/en-vee/alog"

	"github.com/en-vee/aconf"

	"github.com/en-vee/axlrate/core/server"
	"github.com/en-vee/axlrate/service/provisioning"
)

var (
	configFileName = flag.String("config-file-name", "axlrate.conf", "axlrate configuration file. default is axlrate.conf")
)

const (
	PROVISIONING = "Provisioning"
)

type AxlRateConf struct {
	AxlRate struct {
		Server struct {
			Role    string `hocon:"role"`
			Address string `hocon:"address"`
			Port    int64  `hocon:"port"`
		} `hocon:"server"`
	} `hocon:"axlrate"`
}

func main() {

	alog.Info("Parsing Command-Line Arguments")
	// Parse Command-line arguments
	flag.Parse()

	var configParser = aconf.HoconParser{}

	file, err := os.Open(*configFileName)
	defer file.Close()
	if err != nil {
		alog.Critical("Error opening file : %s", *configFileName)
		alog.Critical("Error : %v", err)
		os.Exit(1)
	}

	sysConf := AxlRateConf{}

	if err := configParser.Parse(file, &sysConf); err != nil {

		alog.Critical("Error parsing config file : %s", *configFileName)
		alog.Critical("Error : %v", err)
		os.Exit(1)
	}

	//alog.Info("sysConf=%v", sysConf)
	serverRole := sysConf.AxlRate.Server.Role

	// Setup Signal Handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Based on the role in the config file, launch the appropriate server

	var srv server.Server

	switch serverRole {
	case PROVISIONING:
		srv = &provisioning.Server{NetworkComponent: server.NetworkComponent{NetworkType: server.TcpAddressType, Address: sysConf.AxlRate.Server.Address, PortNumber: sysConf.AxlRate.Server.Port}}
	default:
		alog.Critical("Unknown role : %s", serverRole)
		os.Exit(1)
	}

	errChan := server.Launch(srv)

	for {
		select {
		case sig := <-sigChan:
			alog.Info("Received Signal %v", sig)
		case err := <-errChan:
			alog.Error("Provisioning Server Error : %v", err)
		case <-time.After(time.Second):
			alog.Trace("No error reported by provisioning server in last 1 second")
		}
	}

	alog.Info("Terminating axlrate")
}
