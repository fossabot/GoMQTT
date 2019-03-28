package mqtt

import (
	"context"
	"net"
	"sync"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

type client struct {
	typ        int
	mu         sync.Mutex
	broker     *Broker
	conn       net.Conn
	info       info
	route      route
	status     int
	ctx        context.Context
	cancelFunc context.CancelFunc
	session    *sessions.Session
	subMap     map[string]*subscription
	topicsMgr  *topics.Manager
	subs       []interface{}
	qoss       []byte
	rmsgs      []*packets.PublishPacket
}
