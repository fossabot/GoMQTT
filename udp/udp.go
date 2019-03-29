package udp

import (
	"errors"
	"log"
	"net"
	"sync"
)

var tIndex topicNames
var clients Clients

func validateClientId(clientid []byte) (string, error) {
	if len(clientid) == 0 {
		return "", errors.New("zero-length client id not allowed")
	}
	if len(clientid) > 23 {
		return "", errors.New("client id longer than 23 characters")
	}
	return string(clientid), nil
}

func ListenUDP(addr string) {
	tIndex = topicNames{
		sync.RWMutex{},
		make(map[uint16]string),
		0,
	}
	clients = Clients{
		sync.RWMutex{},
		make(map[string]SNClient),
	}

	address, err := net.ResolveUDPAddr("udp", addr)
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
		go ProcessPacket(n, buf, udpconn, remote)
	}
}
