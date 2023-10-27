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

func main() {
	var iface = flag.String("i", "eth0", "Interface to read packets from")
	var dstIp = flag.String("dst", "127.0.0.1", "dest ip addr")
	var dstPortStr = flag.String("p", "22", "dest port")
	flag.Parse()

	dstPort, _ := strconv.Atoi(*dstPortStr)

	ethernet := &layers.Ethernet{
		BaseLayer:    layers.BaseLayer{},
		SrcMAC:       parseMac("00:00:00:00:00:00"),
		DstMAC:       parseMac("00:00:00:00:00:00"),
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		Flags:    layers.IPv4DontFragment,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		SrcIP:    net.ParseIP("127.0.0.1"),
		DstIP:    net.ParseIP(*dstIp),
	}
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(35000),
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

	// SYNパケットを送信
	handle.WritePacketData(packetbuf.Bytes())
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		//tcpLayer := packet.Layer(layers.LayerTypeTCP)
		//synack := tcpLayer.(*layers.TCP)
		fmt.Printf("recv packet is %+v\n", packet)
	}
}
