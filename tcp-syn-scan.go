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
	"strings"
)

type nwDevice struct {
	macAddr  net.HardwareAddr
	ipv4Addr net.IP
}

func parseMac(macaddr string) net.HardwareAddr {
	parsedMac, _ := net.ParseMAC(macaddr)
	return parsedMac
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

func main() {
	var iface = flag.String("i", "eth0", "Interface to read packets from")
	var dstIp = flag.String("dst", "127.0.0.1", "dest ip addr")
	var dstMac = flag.String("dstmac", "00:00:00:00:00:00", "dest macc addr")
	var dstPortStr = flag.String("p", "22", "dest port")
	flag.Parse()

	var srcPort = 20
	dstPort, _ := strconv.Atoi(*dstPortStr)
	netif := getInterface(*iface)

	ethernet := &layers.Ethernet{
		BaseLayer:    layers.BaseLayer{},
		SrcMAC:       netif.macAddr,
		DstMAC:       parseMac(*dstMac),
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		Flags:    layers.IPv4DontFragment,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		SrcIP:    netif.ipv4Addr,
		DstIP:    net.ParseIP(*dstIp),
	}
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
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
			break
		} else {
			fmt.Printf("TCP %d is close\n", dstPort)
			break
		}
	}
}
