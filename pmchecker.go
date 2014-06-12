package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/oxtoacart/go-igdman/igdman"
)

var (
	internalIP string
)

func main() {
	var err error
	internalIP, err = getFirstNonLoopbackAdapterAddr()
	if err != nil {
		log.Fatalf("Unable to determine internal IP: %s", err)
	}

	upnpIGD, err := igdman.NewUpnpIGD()
	if err != nil {
		log.Printf("UPnP not available: %s", err)
	} else {
		testPorts(upnpIGD, "UPnP")
	}

	natPMPIGD, err := igdman.NewNATPMPIGD()
	if err != nil {
		log.Printf("NAT-PMP not available: %s", err)
	} else {
		testPorts(natPMPIGD, "NAT-PMP")
	}
}

func testPorts(igd igdman.IGD, igdType string) {
	portsToCheck := []int{8443, 443}
	expiration := 1 * time.Second

	for _, externalPort := range portsToCheck {
		igd.RemovePortMapping(igdman.TCP, externalPort)
		err := igd.AddPortMapping(igdman.TCP, internalIP, 15600, externalPort, expiration)
		success := "success"
		if err != nil {
			success = err.Error()
		}
		log.Printf("%s Port %d: %s", igdType, externalPort, success)
	}
}

func getFirstNonLoopbackAdapterAddr() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		ip := net.ParseIP(a)
		if !ip.IsLoopback() {
			return a, nil
		}
	}

	return "", fmt.Errorf("No non-loopback adapter found")
}
