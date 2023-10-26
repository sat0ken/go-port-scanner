package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	targetIP := os.Args[1]
	targetPort := os.Args[2]

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", targetIP, targetPort))
	if err != nil {
		log.Fatalf("TCP %s is close", targetPort)
	}
	defer conn.Close()

	fmt.Printf("TCP %s is open\n", targetPort)
}
