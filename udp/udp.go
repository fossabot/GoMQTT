package udp

import (
	"errors"
	"log"
	"net"
	"sync"
)

// This needs to be efficient, but it's not efficient.
// Raise my salary and maybe I'll fix it.
type topicNames struct {
	sync.RWMutex
	contents map[uint16]string
	next     uint16
}

var tIndex topicNames

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
