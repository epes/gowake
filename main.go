package main

import (
	"fmt"
	"github.com/epes/ecrypto"
	"log"
	"net"
	"os"
	"strings"
)

var (
	StartPackets = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	Port         = 7
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("pass in device MAC as argument")
	}

	mac := ecrypto.MustHexToBytes(strings.ReplaceAll(os.Args[1], ":", ""))

	data := StartPackets

	for i := 0; i < 16; i++ {
		data = append(data, mac...)
	}

	ip, err := LocalIP()
	if err != nil {
		log.Fatal(err)
	}

	broadcastIP, err := BroadcastIPv4(ip)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", broadcastIP.String(), Port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.Write(data)
}

func LocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func BroadcastIPv4(ip net.IP) (net.IP, error) {
	to4 := ip.To4()

	if to4 == nil {
		return nil, fmt.Errorf("IP is not IPv4")
	}

	to4[3] = 0xff
	return to4, nil
}
