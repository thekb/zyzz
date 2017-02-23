package main

import (
	"net"
	"fmt"
	"bytes"
	"encoding/binary"
	"time"
)

func main() {
	serverAddresss, err := net.ResolveUDPAddr("udp", ":10001")
	if err != nil {
		fmt.Println("unable to resolve test_udp address:", err)
		return
	}
	conn, err := net.ListenUDP("udp", serverAddresss)
	if err != nil {
		fmt.Println("unable to listen packet:", err)
		return
	}
	defer conn.Close()
	buffer := make([]byte, 16)
	/*
		payload will contain
		uint16 packet number // little endian
		int64 client time in nano seconds // little endian
		total size of the packet 16 bytes

	 */
	for {
		var packetNumber uint64
		var nanoTime int64
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("error reading from upd:", err)
		} else {
			fmt.Println("received packet from address, bytes:", addr, n)
			if n != 16 {
				fmt.Println("expected 16 bytes of payload")
				continue
			} else {
				buf := bytes.NewReader(buffer[0:8])
				err := binary.Read(buf, binary.LittleEndian, &packetNumber)
				if err != nil {
					fmt.Println("unable to read packet number:", err)
					continue
				}
				fmt.Println("packet number:", packetNumber)

				buf = bytes.NewReader(buffer[8:16])
				err = binary.Read(buf, binary.LittleEndian, &nanoTime)
				if err != nil {
					fmt.Println("unable to read nano time:", err)
					continue
				}
				currentTime := time.Now().UnixNano()
				fmt.Println("nano time:", nanoTime)
				fmt.Println("time taken:", (currentTime - nanoTime)/1000)
			}
		}


	}


}