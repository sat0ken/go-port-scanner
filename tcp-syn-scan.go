package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/routing"
	"log"
	"net"
)

type scanner struct {
	iface        *net.Interface
	dst, gw, src net.IP
	handle       *pcap.Handle
	opts         gopacket.SerializeOptions
	buf          gopacket.SerializeBuffer
}

func newScanner(ip net.IP, router routing.Router) (*scanner, error) {
	s := &scanner{
		dst: ip,
		opts: gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
		buf: gopacket.NewSerializeBuffer(),
	}
	iface, gw, src, err := router.Route(ip)
	if err != nil {
		return nil, err
	}
	log.Printf("scanning ip %v with interface %v, gateway %v, src %v", ip, iface.Name, gw, src)
	s.gw, s.src, s.iface = gw, src, iface

	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	s.handle = handle
	return s, nil
}

func (s *scanner) close() {
	s.handle.Close()
}

func synPacket() {

}

func main() {

}
