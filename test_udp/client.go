package main

import (
	"net"
	"fmt"
	"bytes"
	"encoding/binary"
	"time"
)

func main() {
	serverAddress, err := net.ResolveUDPAddr("udp", "35.154.152.224:10001")
	if err != nil {
		fmt.Println("unable to resolve server address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddress)
	if err != nil {
		fmt.Println("unable to dial test_udp:", err)
		return
	}
	defer conn.Close()

	var packetNumber uint64
	var nanoTime int64
	for {
		packetNumber += 1
		nanoTime = time.Now().UnixNano()
		buffer := new(bytes.Buffer)
		binary.Write(buffer, binary.LittleEndian, packetNumber)
		binary.Write(buffer, binary.LittleEndian, nanoTime)
		conn.Write(buffer.Bytes())
	}
}