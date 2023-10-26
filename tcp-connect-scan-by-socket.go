package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func strToIPByte(str string) (ipbyte [4]byte) {
	splitstr := strings.Split(str, ".")
	for i, v := range splitstr {
		ip, _ := strconv.Atoi(v)
		ipbyte[i] = byte(ip)
	}
	return ipbyte
}

func strToInt(port string) int {
	i, _ := strconv.Atoi(port)
	return i
}

func main() {
	targetIP := os.Args[1]
	targetPort := os.Args[2]

	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		log.Fatalf("create socket err : %v", err)
	}
	err = syscall.Connect(s, &syscall.SockaddrInet4{
		Port: strToInt(targetPort),
		Addr: strToIPByte(targetIP),
	})
	if err != nil {
		log.Fatalf("TCP %s is close", targetPort)
	}
	fmt.Printf("TCP %s is open\n", targetPort)
}
