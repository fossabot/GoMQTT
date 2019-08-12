package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

func ProcessPacket(nbytes int, buffer []byte, con *net.UDPConn, addr *net.UDPAddr) {
	buffer = buffer[:nbytes]
	buf := bytes.NewBuffer(buffer)
	rawmsg, _ := ReadPacket(buf)
	if debug {
		fmt.Println(hex.EncodeToString(buffer))
	}

	switch msg := rawmsg.(type) {
	case *AdvertiseMessage:
		// ADVERTISE must be handled by a broker in future to allow
		// clusterized MQTT-SN clouds: brokers are also (forwarding) clients.
	case *SearchGwMessage:
		// SEARCHGW is useful for searching for new brokers in range of a
		// single network hop. Typically, a broker must NOT broadcast
		// SEARCHGW on more than a single hop.
	case *GwInfoMessage:
		// Each broker must implement GWINFO to supply the automated creation of
		// clusterized MQTT-SN clouds.
	case *ConnectMessage:
		clientid, err := validateClientId(msg.ClientId)
		if err != nil {
			log.Println(err)
			return
		}
		tClient := NewClient(string(clientid), con, addr)
		if msg.Will {

		}
		clients.AddClient(tClient)
		ca := NewMessage(CONNACK).(*ConnackMessage)
		ca.ReturnCode = 0
		if err = tClient.Write(ca); err != nil {
			log.Println(err)
		}
	case *ConnackMessage:
		// CONNACK is a next step of a MQTT-SN cluster system creation. As it was
		// stated earlier, a broker is also a (forwarding) client for other brokers.
	case *WillTopicReqMessage:
		// WILLTOPICREQ lol
	case *WillTopicMessage:
		// WILLTOPIC lol
	case *WillMsgReqMessage:
		// WILLMSGREQ lol
	case *WillMsgMessage:
		// WILLMSQ lol
	case *RegisterMessage:
		topic := string(msg.TopicName)
		var topicid uint16
		if !tIndex.containsTopic(topic) {
			topicid = tIndex.putTopic(topic)
		} else {
			topicid = tIndex.getId(topic)
		}
		tclient := clients.GetClient(addr)
		if tclient == nil {
			log.Println("Received packet from non-existent user!")
			return
		}
		tclient.Register(topicid, topic)
		a := NewMessage(REGACK).(*RegackMessage)
		a.TopicId = topicid
		a.MessageId = msg.MessageId
		a.ReturnCode = 0
		if err := tclient.Write(a); err != nil {
			log.Println(err)
		}
	case *RegackMessage:
		// REGACK may occur on broker level because brokers may also subscribe to
		// (supposedly wildcard) topics on other brokers and forward messages.
	case *PublishMessage:
		topic := tIndex.getTopic(msg.TopicId)
		for _, client := range clients.clients {
			if client.registeredTopics[msg.TopicId] == topic {
				client.Write(msg)
			}
		}
		if msg.Qos > 0 {
			a := NewMessage(PUBACK).(*PubackMessage)
			a.ReturnCode = 0
			a.MessageId = msg.MessageId
			a.TopicId = msg.TopicId
			clients.GetClient(addr).Write(a)
		}
	case *PubackMessage:
		// PUBACK is needed if QoS level between brokers is >0.
	case *PubcompMessage:
		// PUBCOMP is used by MQTT-SN itself to ensure that the message was
		// delivered exactly once.
	case *PubrecMessage:
		// PUBREC is a first message sent in response by a broker on QoS 2
		// to acknowledge the client that the message was received.
	case *PubrelMessage:
		// PUBREL is a next step of MQTT-SN QoS 2 publication acknowledgement
		// process that ensures the publication further, avoiding duplicate
		// publishing.
	case *SubscribeMessage:
		var answer byte
		var topicID uint16
		switch msg.TopicIdType {
		case 0x00, 0x02:
			if !tIndex.containsTopic(string(msg.TopicName)) {
				topicID = tIndex.putTopic(string(msg.TopicName))
			}
		case 0x01:
			if !tIndex.containsId(msg.TopicId) {
				log.Println("requested topic ID not found:", msg.TopicId)
				answer = REJ_INVALID_TID
			}
			topicID = msg.TopicId
		}
		clients.GetClient(addr).Register(topicID, tIndex.getTopic(topicID))
		ack := NewMessage(SUBACK).(*SubackMessage)
		ack.MessageId = msg.MessageId
		ack.Qos = msg.Qos
		ack.ReturnCode = answer
		ack.TopicId = topicID
		tclient := clients.GetClient(addr)
		tclient.Write(ack)
	case *SubackMessage:
		// SUBACK is processed by a broker as well when subscribing to other
		// brokers.
	case *UnsubackMessage:
		// UNSUBACK is processed by a broker as well when subscribing to other
		// brokers.
	case *PingreqMessage:
		a := NewMessage(PINGRESP).(*PingrespMessage)
		clients.GetClient(addr).Write(a)
	case *DisconnectMessage:
		tclient := clients.GetClient(addr)
		msg.Duration = 0
		tclient.Write(msg)
	case *WillTopicUpdateMessage:
		// WILLTOPICUPD lol
	case *WillTopicRespMessage:
		// WILLTOPICRESP lol
	case *WillMsgUpdateMessage:
		// WILLMSGUPD lol
	case *WillMsgRespMessage:
		// WILLMSGRESP lol
	default:
		log.Printf("Unknown Message Type %T\n", msg)
	}
}
