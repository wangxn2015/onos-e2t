// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/wangxn2015/onos-e2t/pkg/manager"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var log = logging.GetLogger()

func main() {
	logging.SetLevel(logging.InfoLevel)
	var serviceModelPlugins arrayFlags
	flag.Var(&serviceModelPlugins, "serviceModel", "names of service model plugins to load (repeated)")
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	sctpPort := flag.Uint("sctpport", 36421, "sctp server port")
	topoEndpoint := flag.String("topoEndpoint", "onos-topo:5150", "onos-topo endpoint address")
	//----------
	e2NodeContainerMode := flag.String("e2NodeContainerMode", "false", "e2NodeContainerMode is false means connecting E2node from outside K8S")
	e2tInterface0IP := flag.String("e2tInterface0IP", "192.168.127.113", "Info for creating E2T, used if e2NodeContainerMode is false")
	e2tInterface0Port := flag.Uint("e2tInterface0Port", 36401, "Info for creating E2T, used if e2NodeContainerMode is false")

	flag.Parse()

	log.Warn("Starting onos-e2t...now")
	cfg := manager.Config{
		CAPath:              *caPath,
		KeyPath:             *keyPath,
		CertPath:            *certPath,
		GRPCPort:            5150,
		E2Port:              int(*sctpPort),
		TopoAddress:         *topoEndpoint,
		ServiceModelPlugins: serviceModelPlugins,
		//----wxn----
		E2NodeContainerMode: *e2NodeContainerMode,
		E2tInterface0IP:     *e2tInterface0IP,
		E2tInterface0Port:   int(*e2tInterface0Port),
	}
	mgr := manager.NewManager(cfg)
	mgr.Run()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	mgr.Close()
}
