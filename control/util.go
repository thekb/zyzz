package control

import (
	"net"
	"fmt"
	"time"
)


func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		fmt.Println("unable to resolve tcp address 0:", err)
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("unabel to listen to tcp address:", err)
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}


func GetCurrentTimeInMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}