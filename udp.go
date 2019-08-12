package main

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
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
		make(map[string]*Client),
	}

	address, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	udpconn, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Fatalln(err)
	}
	go Advertise(180)
	for {
		buf := make([]byte, serv.Config.Buffer)
		n, remote, err := udpconn.ReadFromUDP(buf)
		if err != nil {
			// TODO: Better error processing
			log.Println("Socket error:", err)
			time.Sleep(3 * time.Second)
			continue
		}
		if n < 2 {
			log.Println("Bad data from", remote.String())
			continue
		}
		go ProcessPacket(n, buf, udpconn, remote)
	}
}

// Advertise sends packet nearly every `d` seconds
func Advertise(d uint16) {
	for {
		for _, client := range clients.clients {
			adv := NewMessage(ADVERTISE).(*AdvertiseMessage)
			adv.GatewayId = 0
			adv.Duration = d
			client.Write(adv)
		}
		time.Sleep((time.Duration(d) * time.Second) - (850 * time.Millisecond))
	}
}
