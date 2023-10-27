package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
)

func parseMac(macaddr string) net.HardwareAddr {
	parsedMac, _ := net.ParseMAC(macaddr)
	return parsedMac
}

func main() {
	var iface = flag.String("i", "eth0", "Interface to read packets from")
	flag.Parse()

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
		Protocol: layers.IPProtocolICMPv4,
		SrcIP:    net.ParseIP("127.0.0.1"),
		DstIP:    net.ParseIP("127.0.0.1"),
	}
	icmp := &layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       0,
		Seq:      0,
	}
	packetbuf := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(
		packetbuf,
		gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
		ethernet,
		ip,
		icmp)
	if err != nil {
		log.Fatalf("create packet err : %v", err)
	}

	handle, err := pcap.OpenLive(*iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	// pingを送信
	handle.WritePacketData(packetbuf.Bytes())
	// pingを受信
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
		reply := icmpLayer.(*layers.ICMPv4)
		if reply.TypeCode == layers.ICMPv4TypeEchoReply {
			fmt.Println("recieve echo reply")
			break
		}
	}
}
