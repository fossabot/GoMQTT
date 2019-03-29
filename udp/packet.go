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

func NewMessage(msgType byte) (m Message) {
	switch msgType {
	case ADVERTISE:
		m = &AdvertiseMessage{Header: Header{MessageType: ADVERTISE, Length: 5}}
	case SEARCHGW:
		m = &SearchGwMessage{Header: Header{MessageType: SEARCHGW, Length: 3}}
	case GWINFO:
		m = &GwInfoMessage{Header: Header{MessageType: GWINFO}}
	case CONNECT:
		m = &ConnectMessage{Header: Header{MessageType: CONNECT}, ProtocolId: 0x01}
	case CONNACK:
		m = &ConnackMessage{Header: Header{MessageType: CONNACK, Length: 3}}
	case WILLTOPICREQ:
		m = &WillTopicReqMessage{Header: Header{MessageType: WILLTOPICREQ, Length: 2}}
	case WILLTOPIC:
		m = &WillTopicMessage{Header: Header{MessageType: WILLTOPIC}}
	case WILLMSGREQ:
		m = &WillMsgReqMessage{Header: Header{MessageType: WILLMSGREQ, Length: 2}}
	case WILLMSG:
		m = &WillMsgMessage{Header: Header{MessageType: WILLMSG}}
	case REGISTER:
		m = &RegisterMessage{Header: Header{MessageType: REGISTER}}
	case REGACK:
		m = &RegackMessage{Header: Header{MessageType: REGACK, Length: 7}}
	case PUBLISH:
		m = &PublishMessage{Header: Header{MessageType: PUBLISH}}
	case PUBACK:
		m = &PubackMessage{Header: Header{MessageType: PUBACK, Length: 7}}
	case PUBCOMP:
		m = &PubcompMessage{Header: Header{MessageType: PUBCOMP, Length: 4}}
	case PUBREC:
		m = &PubrecMessage{Header: Header{MessageType: PUBREC, Length: 4}}
	case PUBREL:
		m = &PubrelMessage{Header: Header{MessageType: PUBREL, Length: 4}}
	case SUBSCRIBE:
		m = &SubscribeMessage{Header: Header{MessageType: SUBSCRIBE}}
	case SUBACK:
		m = &SubackMessage{Header: Header{MessageType: SUBACK, Length: 8}}
	case UNSUBSCRIBE:
		m = &UnsubscribeMessage{Header: Header{MessageType: UNSUBSCRIBE}}
	case UNSUBACK:
		m = &UnsubackMessage{Header: Header{MessageType: UNSUBACK, Length: 4}}
	case PINGREQ:
		m = &PingreqMessage{Header: Header{MessageType: PINGREQ}}
	case PINGRESP:
		m = &PingrespMessage{Header: Header{MessageType: PINGRESP, Length: 2}}
	case DISCONNECT:
		m = &DisconnectMessage{Header: Header{MessageType: DISCONNECT}}
	case WILLTOPICUPD:
		m = &WillTopicUpdateMessage{Header: Header{MessageType: WILLTOPICUPD}}
	case WILLTOPICRESP:
		m = &WillTopicRespMessage{Header: Header{MessageType: WILLTOPICRESP, Length: 3}}
	case WILLMSGUPD:
		m = &WillMsgUpdateMessage{Header: Header{MessageType: WILLMSGUPD}}
	case WILLMSGRESP:
		m = &WillMsgRespMessage{Header: Header{MessageType: WILLMSGRESP, Length: 3}}
	}
	return
}

func NewMessageWithHeader(h Header) (m Message) {
	switch h.MessageType {
	case ADVERTISE:
		m = &AdvertiseMessage{Header: h}
	case SEARCHGW:
		m = &SearchGwMessage{Header: h}
	case GWINFO:
		m = &GwInfoMessage{Header: h}
	case CONNECT:
		m = &ConnectMessage{Header: h}
	case CONNACK:
		m = &ConnackMessage{Header: h}
	case WILLTOPICREQ:
		m = &WillTopicReqMessage{Header: h}
	case WILLTOPIC:
		m = &WillTopicMessage{Header: h}
	case WILLMSGREQ:
		m = &WillMsgReqMessage{Header: h}
	case WILLMSG:
		m = &WillMsgMessage{Header: h}
	case REGISTER:
		m = &RegisterMessage{Header: h}
	case REGACK:
		m = &RegackMessage{Header: h}
	case PUBLISH:
		m = &PublishMessage{Header: h}
	case PUBACK:
		m = &PubackMessage{Header: h}
	case PUBCOMP:
		m = &PubcompMessage{Header: h}
	case PUBREC:
		m = &PubrecMessage{Header: h}
	case PUBREL:
		m = &PubrelMessage{Header: h}
	case SUBSCRIBE:
		m = &SubscribeMessage{Header: h}
	case SUBACK:
		m = &SubackMessage{Header: h}
	case UNSUBSCRIBE:
		m = &UnsubscribeMessage{Header: h}
	case UNSUBACK:
		m = &UnsubackMessage{Header: h}
	case PINGREQ:
		m = &PingreqMessage{Header: h}
	case PINGRESP:
		m = &PingrespMessage{Header: h}
	case DISCONNECT:
		m = &DisconnectMessage{Header: h}
	case WILLTOPICUPD:
		m = &WillTopicUpdateMessage{Header: h}
	case WILLTOPICRESP:
		m = &WillTopicRespMessage{Header: h}
	case WILLMSGUPD:
		m = &WillMsgUpdateMessage{Header: h}
	case WILLMSGRESP:
		m = &WillMsgRespMessage{Header: h}
	}
	return
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

type AdvertiseMessage struct {
	Header
	GatewayId byte
	Duration  uint16
}

func (a *AdvertiseMessage) MessageType() byte {
	return ADVERTISE
}

func (a *AdvertiseMessage) Write(w io.Writer) error {
	packet := a.Header.pack()
	packet.WriteByte(ADVERTISE)
	packet.WriteByte(a.GatewayId)
	packet.Write(encodeUint16(a.Duration))
	_, err := packet.WriteTo(w)

	return err
}

func (a *AdvertiseMessage) Unpack(b io.Reader) {
	a.GatewayId = readByte(b)
	a.Duration = readUint16(b)
}

func NewAdvertiseMessage(gwid byte, duration uint16) *AdvertiseMessage {
	return &AdvertiseMessage{
		// TODO
		GatewayId: gwid,
		Duration:  duration,
	}
}

type SearchGwMessage struct {
	Header
	Radius byte
}

func (s *SearchGwMessage) MessageType() byte {
	return SEARCHGW
}

func (s *SearchGwMessage) Write(w io.Writer) error {
	packet := s.Header.pack()
	packet.WriteByte(SEARCHGW)
	packet.WriteByte(s.Radius)
	_, err := packet.WriteTo(w)

	return err
}

func (s *SearchGwMessage) Unpack(b io.Reader) {
	s.Radius = readByte(b)
}

type GwInfoMessage struct {
	Header
	GatewayId      byte
	GatewayAddress []byte
}

func (g *GwInfoMessage) MessageType() byte {
	return GWINFO
}

func (g *GwInfoMessage) Write(w io.Writer) error {
	g.Header.Length = uint16(len(g.GatewayAddress) + 3)
	packet := g.Header.pack()
	packet.WriteByte(GWINFO)
	packet.WriteByte(g.GatewayId)
	packet.Write(g.GatewayAddress)
	_, err := packet.WriteTo(w)

	return err
}

func (g *GwInfoMessage) Unpack(b io.Reader) {
	g.GatewayId = readByte(b)
	if g.Header.Length > 3 {
		b.Read(g.GatewayAddress)
	}
}

type ConnectMessage struct {
	Header
	Will         bool
	CleanSession bool
	ProtocolId   byte
	Duration     uint16
	ClientId     []byte
}

func (c *ConnectMessage) MessageType() byte {
	return CONNECT
}

func (c *ConnectMessage) decodeFlags(b byte) {
	c.Will = (b & WILLFLAG) == WILLFLAG
	c.CleanSession = (b & CLEANSESSION) == CLEANSESSION
}

func (c *ConnectMessage) encodeFlags() byte {
	var b byte
	if c.Will {
		b |= WILLFLAG
	}
	if c.CleanSession {
		b |= CLEANSESSION
	}
	return b
}

func (c *ConnectMessage) Write(w io.Writer) error {
	c.Header.Length = uint16(len(c.ClientId) + 6)
	packet := c.Header.pack()
	packet.WriteByte(CONNECT)
	packet.WriteByte(c.encodeFlags())
	packet.WriteByte(c.ProtocolId)
	packet.Write(encodeUint16(c.Duration))
	packet.Write([]byte(c.ClientId))
	_, err := packet.WriteTo(w)

	return err
}

func (c *ConnectMessage) Unpack(b io.Reader) {
	c.decodeFlags(readByte(b))
	c.ProtocolId = readByte(b)
	c.Duration = readUint16(b)
	c.ClientId = make([]byte, c.Header.Length-6)
	b.Read(c.ClientId)
}

type ConnackMessage struct {
	Header
	ReturnCode byte
}

func (c *ConnackMessage) MessageType() byte {
	return CONNACK
}

func (c *ConnackMessage) Write(w io.Writer) error {
	packet := c.Header.pack()
	packet.WriteByte(CONNACK)
	packet.WriteByte(c.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (c *ConnackMessage) Unpack(b io.Reader) {
	c.ReturnCode = readByte(b)
}

type WillTopicReqMessage struct {
	Header
}

func (wt *WillTopicReqMessage) MessageType() byte {
	return WILLTOPICREQ
}

func (wt *WillTopicReqMessage) Write(w io.Writer) error {
	packet := wt.Header.pack()
	packet.WriteByte(WILLTOPICREQ)
	_, err := packet.WriteTo(w)

	return err
}

func (wt *WillTopicReqMessage) Unpack(b io.Reader) {

}

type WillTopicMessage struct {
	Header
	Qos       byte
	Retain    bool
	WillTopic []byte
}

func (wt *WillTopicMessage) MessageType() byte {
	return wt.Header.MessageType
}

func (wt *WillTopicMessage) encodeFlags() byte {
	var b byte

	b |= (wt.Qos << 5) & QOSBITS
	if wt.Retain {
		b |= RETAINFLAG
	}
	return b
}

func (wt *WillTopicMessage) decodeFlags(b byte) {
	wt.Qos = (b & QOSBITS) >> 5
	wt.Retain = (b & RETAINFLAG) == RETAINFLAG
}

func (wt *WillTopicMessage) Write(w io.Writer) error {
	if len(wt.WillTopic) == 0 {
		wt.Header.Length = 2
	} else {
		wt.Header.Length = uint16(len(wt.WillTopic) + 3)
	}
	packet := wt.Header.pack()
	packet.WriteByte(wt.Header.MessageType)
	if wt.Header.Length > 2 {
		packet.WriteByte(wt.encodeFlags())
		packet.Write(wt.WillTopic)
	}
	_, err := packet.WriteTo(w)

	return err
}

func (wt *WillTopicMessage) Unpack(b io.Reader) {
	if wt.Header.Length > 2 {
		wt.decodeFlags(readByte(b))
		b.Read(wt.WillTopic)
	}
}

type WillMsgReqMessage struct {
	Header
}

func (wm *WillMsgReqMessage) MessageType() byte {
	return WILLMSGREQ
}

func (wm *WillMsgReqMessage) Write(w io.Writer) error {
	packet := wm.Header.pack()
	packet.WriteByte(wm.Header.MessageType)
	_, err := packet.WriteTo(w)

	return err
}

func (wm *WillMsgReqMessage) Unpack(b io.Reader) {

}

type WillMsgMessage struct {
	Header
	WillMsg []byte
}

func (wm *WillMsgMessage) MessageType() byte {
	return WILLMSG
}

func (wm *WillMsgMessage) Write(w io.Writer) error {
	wm.Header.Length = uint16(len(wm.WillMsg) + 2)
	packet := wm.Header.pack()
	packet.WriteByte(WILLMSG)
	packet.Write(wm.WillMsg)
	_, err := packet.WriteTo(w)

	return err
}

func (wm *WillMsgMessage) Unpack(b io.Reader) {
	b.Read(wm.WillMsg)
}

type RegisterMessage struct {
	Header
	TopicId   uint16
	MessageId uint16
	TopicName []byte
}

func NewRegisterMessage(TopicId, MessageId uint16, TopicName []byte) *RegisterMessage {
	return &RegisterMessage{
		TopicId:   TopicId,
		MessageId: MessageId,
		TopicName: TopicName,
	}
}

func (r *RegisterMessage) MessageType() byte {
	return REGISTER
}

func (r *RegisterMessage) Write(w io.Writer) error {
	r.Header.Length = uint16(len(r.TopicName) + 6)
	packet := r.Header.pack()
	packet.WriteByte(REGISTER)
	packet.Write(encodeUint16(r.TopicId))
	packet.Write(encodeUint16(r.MessageId))
	packet.Write(r.TopicName)
	_, err := packet.WriteTo(w)

	return err
}

func (r *RegisterMessage) Unpack(b io.Reader) {
	r.TopicId = readUint16(b)
	r.MessageId = readUint16(b)
	r.TopicName = make([]byte, r.Header.Length-6)
	b.Read(r.TopicName)
}

type RegackMessage struct {
	Header
	TopicId    uint16
	MessageId  uint16
	ReturnCode byte
}

func (r *RegackMessage) MessageType() byte {
	return REGACK
}

func (r *RegackMessage) Write(w io.Writer) error {
	packet := r.Header.pack()
	packet.WriteByte(REGACK)
	packet.Write(encodeUint16(r.TopicId))
	packet.Write(encodeUint16(r.MessageId))
	packet.WriteByte(r.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (r *RegackMessage) Unpack(b io.Reader) {
	r.TopicId = readUint16(b)
	r.MessageId = readUint16(b)
	r.ReturnCode = readByte(b)
}

type PublishMessage struct {
	Header
	Dup         bool
	Retain      bool
	Qos         byte
	TopicIdType byte
	TopicId     uint16
	MessageId   uint16
	Data        []byte
}

func NewPublishMessage(TopicId uint16, TopicIdType byte, Data []byte, Qos byte, MessageId uint16, Retain bool, Dup bool) *PublishMessage {
	return &PublishMessage{
		TopicId:     TopicId,
		TopicIdType: TopicIdType,
		Data:        Data,
		Qos:         Qos,
		MessageId:   MessageId,
		Retain:      Retain,
		Dup:         Dup,
	}
}

func (p *PublishMessage) MessageType() byte {
	return PUBLISH
}

func (p *PublishMessage) encodeFlags() byte {
	var b byte
	if p.Dup {
		b |= DUPFLAG
	}
	b |= (p.Qos << 5) & QOSBITS
	if p.Retain {
		b |= RETAINFLAG
	}
	b |= p.TopicIdType & TOPICIDTYPE
	return b
}

func (p *PublishMessage) decodeFlags(b byte) {
	p.Dup = (b & DUPFLAG) == DUPFLAG
	p.Qos = (b & QOSBITS) >> 5
	p.Retain = (b & RETAINFLAG) == RETAINFLAG
	p.TopicIdType = b & TOPICIDTYPE
}

func (p *PublishMessage) Write(w io.Writer) error {
	p.Header.Length = uint16(len(p.Data) + 7)
	packet := p.Header.pack()
	packet.WriteByte(PUBLISH)
	packet.WriteByte(p.encodeFlags())
	packet.Write(encodeUint16(p.TopicId))
	packet.Write(encodeUint16(p.MessageId))
	packet.Write(p.Data)
	_, err := packet.WriteTo(w)

	return err
}

func (p *PublishMessage) Unpack(b io.Reader) {
	p.decodeFlags(readByte(b))
	p.TopicId = readUint16(b)
	p.MessageId = readUint16(b)
	p.Data = make([]byte, p.Header.Length-7)
	b.Read(p.Data)
}

type PubackMessage struct {
	Header
	TopicId    uint16
	MessageId  uint16
	ReturnCode byte
}

func (p *PubackMessage) MessageType() byte {
	return PUBACK
}

func (p *PubackMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PUBACK)
	packet.Write(encodeUint16(p.TopicId))
	packet.Write(encodeUint16(p.MessageId))
	packet.WriteByte(p.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (p *PubackMessage) Unpack(b io.Reader) {
	p.TopicId = readUint16(b)
	p.MessageId = readUint16(b)
	p.ReturnCode = readByte(b)
}

type PubcompMessage struct {
	Header
	MessageId uint16
}

func (p *PubcompMessage) MessageType() byte {
	return PUBCOMP
}

func (p *PubcompMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PUBCOMP)
	packet.Write(encodeUint16(p.MessageId))
	_, err := packet.WriteTo(w)

	return err
}

func (p *PubcompMessage) Unpack(b io.Reader) {
	p.MessageId = readUint16(b)
}

type PubrecMessage struct {
	Header
	MessageId uint16
}

func (p *PubrecMessage) MessageType() byte {
	return PUBREC
}

func (p *PubrecMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PUBREC)
	packet.Write(encodeUint16(p.MessageId))
	_, err := packet.WriteTo(w)

	return err
}

func (p *PubrecMessage) Unpack(b io.Reader) {
	p.MessageId = readUint16(b)
}

type PubrelMessage struct {
	Header
	MessageId uint16
}

func (p *PubrelMessage) MessageType() byte {
	return PUBREL
}

func (p *PubrelMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PUBREL)
	packet.Write(encodeUint16(p.MessageId))
	_, err := packet.WriteTo(w)

	return err
}

func (p *PubrelMessage) Unpack(b io.Reader) {
	p.MessageId = readUint16(b)
}

type SubscribeMessage struct {
	Header
	Dup         bool
	Qos         byte
	TopicIdType byte
	MessageId   uint16
	TopicId     uint16
	TopicName   []byte
}

func (s *SubscribeMessage) MessageType() byte {
	return SUBSCRIBE
}

func (s *SubscribeMessage) encodeFlags() byte {
	var b byte
	if s.Dup {
		b |= DUPFLAG
	}
	b |= (s.Qos << 5) & QOSBITS
	b |= s.TopicIdType & TOPICIDTYPE
	return b
}

func (s *SubscribeMessage) decodeFlags(b byte) {
	s.Dup = (b & DUPFLAG) == DUPFLAG
	s.Qos = (b & QOSBITS) >> 5
	s.TopicIdType = b & TOPICIDTYPE
}

func (s *SubscribeMessage) Write(w io.Writer) error {
	switch s.TopicIdType {
	case 0x00, 0x02:
		s.Header.Length = uint16(len(s.TopicName) + 5)
	case 0x01:
		s.Header.Length = 7
	}
	packet := s.Header.pack()
	packet.WriteByte(SUBSCRIBE)
	packet.WriteByte(s.encodeFlags())
	packet.Write(encodeUint16(s.MessageId))
	switch s.TopicIdType {
	case 0x00, 0x02:
		packet.Write(s.TopicName)
	case 0x01:
		packet.Write(encodeUint16(s.TopicId))
	}
	_, err := packet.WriteTo(w)

	return err
}

func (s *SubscribeMessage) Unpack(b io.Reader) {
	s.decodeFlags(readByte(b))
	s.MessageId = readUint16(b)
	switch s.TopicIdType {
	case 0x00, 0x02:
		s.TopicName = make([]byte, s.Header.Length-5)
		b.Read(s.TopicName)
	case 0x01:
		s.TopicId = readUint16(b)
	}
}

type SubackMessage struct {
	Header
	Qos        byte
	ReturnCode byte
	TopicId    uint16
	MessageId  uint16
}

func NewSubackMessage(TopicId uint16, MessageId uint16, Qos byte, rc byte) *SubackMessage {
	return &SubackMessage{
		Qos:        Qos,
		ReturnCode: rc,
		TopicId:    TopicId,
		MessageId:  MessageId,
	}
}

func (s *SubackMessage) MessageType() byte {
	return SUBACK
}

func (s *SubackMessage) encodeFlags() byte {
	var b byte
	b |= (s.Qos << 5) & QOSBITS
	return b
}

func (s *SubackMessage) decodeFlags(b byte) {
	s.Qos = (b & QOSBITS) >> 5
}

func (s *SubackMessage) Write(w io.Writer) error {
	packet := s.Header.pack()
	packet.WriteByte(SUBACK)
	packet.WriteByte(s.encodeFlags())
	packet.Write(encodeUint16(s.TopicId))
	packet.Write(encodeUint16(s.MessageId))
	packet.WriteByte(s.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (s *SubackMessage) Unpack(b io.Reader) {
	s.decodeFlags(readByte(b))
	s.TopicId = readUint16(b)
	s.MessageId = readUint16(b)
	s.ReturnCode = readByte(b)
}

type UnsubscribeMessage struct {
	Header
	TopicIdType byte
	MessageId   uint16
	TopicId     uint16
	TopicName   []byte
}

func (u *UnsubscribeMessage) MessageType() byte {
	return UNSUBSCRIBE
}

func (s *UnsubscribeMessage) encodeFlags() byte {
	var b byte
	b |= s.TopicIdType & TOPICIDTYPE
	return b
}

func (s *UnsubscribeMessage) decodeFlags(b byte) {
	s.TopicIdType = b & TOPICIDTYPE
}

func (u *UnsubscribeMessage) Write(w io.Writer) error {
	switch u.TopicIdType {
	case 0x00, 0x02:
		u.Header.Length = uint16(len(u.TopicName) + 5)
	case 0x01:
		u.Header.Length = 7
	}
	packet := u.Header.pack()
	packet.WriteByte(UNSUBSCRIBE)
	packet.WriteByte(u.encodeFlags())
	packet.Write(encodeUint16(u.MessageId))
	switch u.TopicIdType {
	case 0x00, 0x02:
		packet.Write(u.TopicName)
	case 0x01:
		packet.Write(encodeUint16(u.TopicId))
	}
	_, err := packet.WriteTo(w)

	return err
}

func (u *UnsubscribeMessage) Unpack(b io.Reader) {
	u.decodeFlags(readByte(b))
	u.MessageId = readUint16(b)
	switch u.TopicIdType {
	case 0x00, 0x02:
		b.Read(u.TopicName)
	case 0x01:
		u.TopicId = readUint16(b)
	}
}

type UnsubackMessage struct {
	Header
	MessageId uint16
}

func (u *UnsubackMessage) MessageType() byte {
	return UNSUBACK
}

func (u *UnsubackMessage) Write(w io.Writer) error {
	packet := u.Header.pack()
	packet.WriteByte(UNSUBACK)
	packet.Write(encodeUint16(u.MessageId))
	_, err := packet.WriteTo(w)

	return err
}

func (u *UnsubackMessage) Unpack(b io.Reader) {
	u.MessageId = readUint16(b)
}

type PingreqMessage struct {
	Header
	ClientId []byte
}

func (p *PingreqMessage) MessageType() byte {
	return PINGREQ
}

func (p *PingreqMessage) Write(w io.Writer) error {
	p.Header.Length = uint16(len(p.ClientId) + 2)
	packet := p.Header.pack()
	packet.WriteByte(PINGREQ)
	if len(p.ClientId) > 0 {
		packet.Write(p.ClientId)
	}
	_, err := packet.WriteTo(w)

	return err
}

func (p *PingreqMessage) Unpack(b io.Reader) {
	if p.Header.Length > 2 {
		b.Read(p.ClientId)
	}
}

type PingrespMessage struct {
	Header
}

func (p *PingrespMessage) MessageType() byte {
	return PINGRESP
}

func (p *PingrespMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PINGRESP)
	_, err := packet.WriteTo(w)

	return err
}

func (p *PingrespMessage) Unpack(b io.Reader) {

}

type DisconnectMessage struct {
	Header
	Duration uint16
}

func (d *DisconnectMessage) MessageType() byte {
	return DISCONNECT
}

func (d *DisconnectMessage) Write(w io.Writer) error {
	var packet bytes.Buffer

	if d.Duration == 0 {
		d.Header.Length = 2
		packet = d.Header.pack()
		packet.WriteByte(DISCONNECT)
	} else {
		d.Header.Length = 4
		packet = d.Header.pack()
		packet.WriteByte(DISCONNECT)
		packet.Write(encodeUint16(d.Duration))
	}
	_, err := packet.WriteTo(w)

	return err
}

func (d *DisconnectMessage) Unpack(b io.Reader) {
	if d.Header.Length == 4 {
		d.Duration = readUint16(b)
	}
}

type WillTopicUpdateMessage struct {
	Header
	Qos       byte
	Retain    bool
	WillTopic []byte
}

func (wt *WillTopicUpdateMessage) MessageType() byte {
	return WILLTOPICUPD
}

func (wt *WillTopicUpdateMessage) encodeFlags() byte {
	var b byte
	b |= (wt.Qos << 5) & QOSBITS
	if wt.Retain {
		b |= RETAINFLAG
	}
	return b
}

func (wt *WillTopicUpdateMessage) decodeFlags(b byte) {
	wt.Qos = (b & QOSBITS) >> 5
	wt.Retain = (b & RETAINFLAG) == RETAINFLAG
}

func (wt *WillTopicUpdateMessage) Write(w io.Writer) error {
	wt.Header.Length = uint16(len(wt.WillTopic) + 3)
	packet := wt.Header.pack()
	packet.WriteByte(WILLTOPICUPD)
	packet.WriteByte(wt.encodeFlags())
	packet.Write(wt.WillTopic)
	_, err := packet.WriteTo(w)

	return err
}

func (wt *WillTopicUpdateMessage) Unpack(b io.Reader) {
	wt.decodeFlags(readByte(b))
	b.Read(wt.WillTopic)
}

type WillTopicRespMessage struct {
	Header
	ReturnCode byte
}

func (wt *WillTopicRespMessage) MessageType() byte {
	return WILLTOPICRESP
}

func (wt *WillTopicRespMessage) Write(w io.Writer) error {
	packet := wt.Header.pack()
	packet.WriteByte(WILLTOPICRESP)
	packet.WriteByte(wt.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (wt *WillTopicRespMessage) Unpack(b io.Reader) {
	wt.ReturnCode = readByte(b)
}

type WillMsgUpdateMessage struct {
	Header
	WillMsg []byte
}

func (wm *WillMsgUpdateMessage) MessageType() byte {
	return WILLMSGUPD
}

func (wm *WillMsgUpdateMessage) Write(w io.Writer) error {
	wm.Header.Length = uint16(len(wm.WillMsg) + 2)
	packet := wm.Header.pack()
	packet.WriteByte(WILLMSGUPD)
	packet.Write(wm.WillMsg)
	_, err := packet.WriteTo(w)

	return err
}

func (wm *WillMsgUpdateMessage) Unpack(b io.Reader) {
	b.Read(wm.WillMsg)
}

type WillMsgRespMessage struct {
	Header
	ReturnCode byte
}

func (wm *WillMsgRespMessage) MessageType() byte {
	return WILLMSGRESP
}

func (wm *WillMsgRespMessage) Write(w io.Writer) error {
	packet := wm.Header.pack()
	packet.WriteByte(WILLMSGRESP)
	packet.WriteByte(wm.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (wm *WillMsgRespMessage) Unpack(b io.Reader) {
	wm.ReturnCode = readByte(b)
}
