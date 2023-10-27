package main

import (
	"flag"
	"github.com/google/gopacket/dumpcommand"
	"github.com/google/gopacket/pcap"
	"log"
)

func main() {
	var iface = flag.String("i", "eth0", "Interface to read packets from")
	flag.Parse()
	handle, err := pcap.OpenLive(*iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	dumpcommand.Run(handle)
}
