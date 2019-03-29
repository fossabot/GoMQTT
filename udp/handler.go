package udp

import (
	"bytes"
	"log"
	"net"
)

func ProcessPacket(nbytes int, buffer []byte, con *net.UDPConn, addr *net.UDPAddr) {
	buf := bytes.NewBuffer(buffer)
	rawmsg, _ := ReadPacket(buf)

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
		// REGISTER lol
	case *RegackMessage:
		// REGACK lol
	case *PublishMessage:
		// PUBLISH omegalul
	case *PubackMessage:
		// PUBACK lol
	case *PubcompMessage:
		// PUBCOMP lol
	case *PubrecMessage:
		// PUBREC lol
	case *PubrelMessage:
		// PUBREL lol
	case *SubscribeMessage:
		// SUBSCRIBE lol
	case *SubackMessage:
		// SUBACK lol
	case *UnsubackMessage:
		// UNSUBACK lol
	case *PingreqMessage:
		// PINGREQ lol
	case *DisconnectMessage:
		// DISCONNECT lol
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
