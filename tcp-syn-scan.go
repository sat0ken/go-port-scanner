package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"strconv"
)

func parseMac(macaddr string) net.HardwareAddr {
	parsedMac, _ := net.ParseMAC(macaddr)
	return parsedMac
}

func getInterface(ifname string) net.Interface {
	netifs, _ := net.Interfaces()
	for _, netif := range netifs {
		if netif.Name == ifname {
			return netif
		}
	}
	return net.Interface{}
}

func main() {
	var iface = flag.String("i", "eth0", "Interface to read packets from")
	var dstIp = flag.String("dst", "127.0.0.1", "dest ip addr")
	var dstPortStr = flag.String("p", "22", "dest port")
	flag.Parse()

	dstPort, _ := strconv.Atoi(*dstPortStr)
	netif := getInterface(*iface)

	ethernet := &layers.Ethernet{
		BaseLayer:    layers.BaseLayer{},
		SrcMAC:       netif.HardwareAddr,
		DstMAC:       parseMac("FC:4A:E9:DE:AD:05"),
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		Flags:    layers.IPv4DontFragment,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		SrcIP:    net.ParseIP("192.168.0.12"),
		DstIP:    net.ParseIP(*dstIp),
	}
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(58742),
		DstPort: layers.TCPPort(dstPort),
		SYN:     true,
	}
	packetbuf := gopacket.NewSerializeBuffer()
	tcp.SetNetworkLayerForChecksum(ip)
	err := gopacket.SerializeLayers(
		packetbuf,
		gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
		ethernet,
		ip,
		tcp)
	if err != nil {
		log.Fatalf("create packet err : %v", err)
	}

	handle, err := pcap.OpenLive(*iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// SYNパケットを送信
	handle.WritePacketData(packetbuf.Bytes())
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	// SYNパケットを受信
	for packet := range packetSource.Packets() {
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		synackip := ipLayer.(*layers.IPv4)
		synack := tcpLayer.(*layers.TCP)
		if synackip.SrcIP.Equal(ip.DstIP) && synack.ACK {
			fmt.Printf("TCP %d is open\n", dstPort)
		} else {
			fmt.Printf("TCP %d is close\n", dstPort)
		}
		break
	}
}
