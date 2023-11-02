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
	// コマンド引数を受け取る
	var iface = flag.String("i", "eth0", "Interface to read packets from")
	var dstIp = flag.String("dst", "192.168.1.3", "dest ip addr")
	var dstMac = flag.String("dstmac", "", "dest ip addr")
	flag.Parse()

	// 送信元のMACアドレスとIPアドレスを取得する
	netif := getInterface(*iface)

	// Ethernetのパケットを作成する
	ethernet := &layers.Ethernet{
		SrcMAC:       netif.macAddr,
		DstMAC:       parseMac(*dstMac),
		EthernetType: layers.EthernetTypeIPv4,
	}
	// IPヘッダを作成する
	ip := &layers.IPv4{
		Version:  4,
		Flags:    layers.IPv4DontFragment,
		TTL:      64,
		Protocol: layers.IPProtocolICMPv4,
		SrcIP:    netif.ipv4Addr,
		DstIP:    net.ParseIP(*dstIp),
	}
	// ICMP Echo Requestのパケットを作成する
	icmp := &layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       0,
		Seq:      0,
	}
	packetbuf := gopacket.NewSerializeBuffer()
	// パケットを作成する
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
			fmt.Printf("recieve echo reply %+v\n", reply)
			break
		}
	}
}
