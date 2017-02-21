package main

import (
	"net"
	"fmt"
	"bytes"
	"encoding/binary"
	"time"
)

func main() {
	serverAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	if err != nil {
		fmt.Println("unable to resolve server address:", err)
		return
	}

	localAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("unable to resolve local address:", err)
		return
	}

	conn, err := net.DialUDP("udp", localAddress, serverAddress)
	if err != nil {
		fmt.Println("unable to dial udp:", err)
		return
	}
	defer conn.Close()

	var packetNumber uint16
	var nanoTime int64
	for i := 0; i < 5; i++ {
		packetNumber += 1
		nanoTime = time.Now().UnixNano()
		buffer := new(bytes.Buffer)
		binary.Write(buffer, binary.LittleEndian, packetNumber)
		binary.Write(buffer, binary.LittleEndian, nanoTime)
		conn.Write(buffer.Bytes())
	}
}