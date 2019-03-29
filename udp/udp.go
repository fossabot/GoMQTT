package udp

import (
	"log"
	"net"
)

func ListenUDP() {
	address, err := net.ResolveUDPAddr("udp", serv.Config.MQTTSNAddress)
	if err != nil {
		log.Fatalln(err)
	}
	udpconn, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		buf := make([]byte, 1024)
		n, remote, err := udpconn.ReadFromUDP(buf)
		if err != nil {
			// TODO: better error processing
			log.Fatalln(err)
		}

	}
}
