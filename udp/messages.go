package udp

// This is going to be an ennormously large file

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Flags
const (
	TOPICIDTYPE  = 0x03
	CLEANSESSION = 0x04
	WILLFLAG     = 0x08
	RETAINFLAG   = 0x10
	QOSBITS      = 0x60
	DUPFLAG      = 0x80
)

// Errors
const (
	ACCEPTED         = 0x00
	REJ_CONGESTION   = 0x01
	REJ_INVALID_TID  = 0x02
	REJ_NOT_SUPORTED = 0x03
)

// Message Types
const (
	ADVERTISE     = 0x00
	SEARCHGW      = 0x01
	GWINFO        = 0x02
	CONNECT       = 0x04
	CONNACK       = 0x05
	WILLTOPICREQ  = 0x06
	WILLTOPIC     = 0x07
	WILLMSGREQ    = 0x08
	WILLMSG       = 0x09
	REGISTER      = 0x0A
	REGACK        = 0x0B
	PUBLISH       = 0x0C
	PUBACK        = 0x0D
	PUBCOMP       = 0x0E
	PUBREC        = 0x0F
	PUBREL        = 0x10
	SUBSCRIBE     = 0x12
	SUBACK        = 0x13
	UNSUBSCRIBE   = 0x14
	UNSUBACK      = 0x15
	PINGREQ       = 0x16
	PINGRESP      = 0x17
	DISCONNECT    = 0x18
	WILLTOPICUPD  = 0x1A
	WILLTOPICRESP = 0x1B
	WILLMSGUPD    = 0x1C
	WILLMSGRESP   = 0x1D
	// 0x03 is reserved
	// 0x11 is reserved
	// 0x19 is reserved
	// 0x1E - 0xFD is reserved
	// 0xFE - Encapsulated message
	// 0xFF is reserved
)

var MessageNames = map[byte]string{
	ADVERTISE:     "ADVERTISE",
	SEARCHGW:      "SEARCHGW",
	GWINFO:        "GWINFO",
	CONNECT:       "CONNECT",
	CONNACK:       "CONNACK",
	WILLTOPICREQ:  "WILLTOPICREQ",
	WILLTOPIC:     "WILLTOPIC",
	WILLMSGREQ:    "WILLMSGREQ",
	WILLMSG:       "WILLMSG",
	REGISTER:      "REGISTER",
	REGACK:        "REGACK",
	PUBLISH:       "PUBLISH",
	PUBACK:        "PUBACK",
	PUBCOMP:       "PUBCOMP",
	PUBREC:        "PUBREC",
	PUBREL:        "PUBREL",
	SUBSCRIBE:     "SUBSCRIBE",
	SUBACK:        "SUBACK",
	UNSUBSCRIBE:   "UNSUBSCRIBE",
	UNSUBACK:      "UNSUBACK",
	PINGREQ:       "PINGREQ",
	PINGRESP:      "PINGRESP",
	DISCONNECT:    "DISCONNECT",
	WILLTOPICUPD:  "WILLTOPICUPD",
	WILLTOPICRESP: "WILLTOPICRESP",
	WILLMSGUPD:    "WILLMSGUPD",
	WILLMSGRESP:   "WILLMSGRESP",
}

type Header struct {
	Length      uint16
	MessageType byte
}

func (h *Header) unpack(b io.Reader) {
	lengthCheck := readByte(b)
	if lengthCheck == 0x01 {
		h.Length = readUint16(b)
	} else {
		h.Length = uint16(lengthCheck)
	}
	h.MessageType = readByte(b)
}

func (h *Header) pack() bytes.Buffer {
	var header bytes.Buffer
	if h.Length > 256 {
		h.Length += 2
		header.WriteByte(0x01)
		header.Write(encodeUint16(h.Length))
	} else {
		header.WriteByte(byte(h.Length))
	}
	return header
}

type Message interface {
	MessageType() byte
	Write(io.Writer) error
	Unpack(io.Reader)
}

func ReadPacket(r io.Reader) (m Message, err error) {
	var h Header
	packet := make([]byte, 1500)
	r.Read(packet)
	packetBuf := bytes.NewBuffer(packet)
	h.unpack(packetBuf)
	m = NewMessageWithHeader(h)
	if m == nil {
		return nil, errors.New("Bad data from client")
	}
	m.Unpack(packetBuf)
	return m, nil
}

func readByte(b io.Reader) byte {
	num := make([]byte, 1)
	b.Read(num)
	return num[0]
}

func readUint16(b io.Reader) uint16 {
	num := make([]byte, 2)
	b.Read(num)
	return binary.BigEndian.Uint16(num)
}

func encodeUint16(num uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, num)
	return bytes
}
