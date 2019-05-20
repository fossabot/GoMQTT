package main

import (
	"bytes"
	"log"
	"net"
)

func ProcessPacket(nbytes int, buffer []byte, con *net.UDPConn, addr *net.UDPAddr) {
	buf := bytes.NewBuffer(buffer)
	rawmsg, _ := ReadPacket(buf)
	debug(rawmsg)

	switch msg := rawmsg.(type) {
	case *AdvertiseMessage:
		// ADVERTISE lol
	case *SearchGwMessage:
		// SEARCHGW lol
	case *GwInfoMessage:
		// GWINFO lol
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
		// CONNACK lol
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
		tclient.Register(topicid, topic)
		a := NewMessage(REGACK).(*RegackMessage)
		a.TopicId = topicid
		a.MessageId = msg.MessageId
		a.ReturnCode = 0
		if err := tclient.Write(a); err != nil {
			log.Println(err)
		}
	case *RegackMessage:
		// REGACK lol
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
		// PUBACK lol
	case *PubcompMessage:
		// PUBCOMP lol
	case *PubrecMessage:
		// PUBREC lol
	case *PubrelMessage:
		// PUBREL lol
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
		// SUBACK lol
	case *UnsubackMessage:
		// UNSUBACK lol
	case *PingreqMessage:
		// PINGREQ lol
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
		log.Println("Unknown Message Type %T\n", msg)
	}
}
