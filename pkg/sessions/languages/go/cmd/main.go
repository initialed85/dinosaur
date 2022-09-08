package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	port = 13337
)

func receiveCallback(conn net.PacketConn, addr net.Addr, data []byte, localAddr *net.UDPAddr) {
	if addr.String() == localAddr.String() {
		return
	}

	fmt.Printf("%v:%v\t%v\n", addr.(*net.UDPAddr).IP.String(), addr.(*net.UDPAddr).Port, string(data))
}

func receiveLoop(conn net.PacketConn, localAddr *net.UDPAddr) {
	buf := make([]byte, 65507)

	for {
		err := conn.SetDeadline(time.Now().Add(time.Second))
		if err != nil {
			log.Fatal(err)
		}

		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		receiveCallback(conn, addr, buf[:n], localAddr)
	}
}

func main() {
	hostname := os.Getenv("HOSTNAME")
	localIP := os.Getenv("LOCAL_IP")
	broadcastIP := os.Getenv("BROADCAST_IP")

	localAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", localIP, port))
	if err != nil {
		log.Fatal(err)
	}

	broadcastAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", broadcastIP, port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenPacket("udp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatal(err)
	}

	go receiveLoop(conn, localAddr)

	data := []byte(fmt.Sprintf("Hello world from Go @ %v", hostname))

	for {
		_, _ = conn.WriteTo(data, broadcastAddr)
		time.Sleep(time.Second)
	}
}
