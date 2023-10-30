package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"strings"
)

type nwDevice struct {
	macAddr  net.HardwareAddr
	ipv4Addr net.IP
}

func getInterface(ifname string) nwDevice {
	netifs, _ := net.Interfaces()
	for _, netif := range netifs {
		if netif.Name == ifname {
			addrs, _ := netif.Addrs()
			for _, addr := range addrs {
				if !strings.Contains(addr.String(), ":") && strings.Contains(addr.String(), ".") {
					ip, _, _ := net.ParseCIDR(addr.String())
					return nwDevice{
						macAddr:  netif.HardwareAddr,
						ipv4Addr: ip,
					}
				}
			}

		}
	}
	return nwDevice{}
}

func parseMac(macaddr string) net.HardwareAddr {
	parsedMac, _ := net.ParseMAC(macaddr)
	return parsedMac
}

func main() {
	var iface = flag.String("i", "eth0", "Interface to read packets from")
	//var dstIp = flag.String("dst", "127.0.0.1", "dest ip addr")
	flag.Parse()

	netif := getInterface(*iface)

	ethernet := &layers.Ethernet{
		SrcMAC:       netif.macAddr,
		DstMAC:       parseMac("FF:FF:FF:FF:FF:FF"),
		EthernetType: layers.EthernetTypeARP,
	}
	arpreq := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   netif.macAddr,
		SourceProtAddress: []byte{192, 168, 0, 56},
		DstHwAddress:      []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    []byte{192, 168, 0, 1},
	}
	packetbuf := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(
		packetbuf,
		gopacket.SerializeOptions{
			FixLengths: true,
		},
		ethernet,
		arpreq)
	if err != nil {
		log.Fatalf("create packet err : %v", err)
	}
	handle, err := pcap.OpenLive(*iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	// ARPリクエストを送信
	handle.WritePacketData(packetbuf.Bytes())
	// ARPリプライを受信
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		ethernetPacket := ethernetLayer.(*layers.Ethernet)
		if ethernetPacket.EthernetType.LayerType() == layers.LayerTypeARP {
			fmt.Printf("packet is %+v\n", packet)
			break
		}
	}
}
