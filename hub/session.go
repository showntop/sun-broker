package hub

import (
	// "github.com/showntop/sun-broker/store"
	"github.com/gomqtt/packet"
)

type Session interface {
	AddSubscription(*packet.Subscription) error
	HoldPacket(packet.Packet) error
	LookupPacket(uint16) (packet.Packet, error)
	PileupMsg(*packet.Message) error
	PublishMsg(*packet.PublishPacket) error
	GetOutChan() chan *packet.Message
}

type ToughSession struct {
	oid         string
	hub         *Hub
	outChan     chan *packet.Message
	holdPackets map[uint16]packet.Packet

	offlineMsgs []string

	will *packet.Message

	currentClient *remoteClient
	clean         bool
}

func NewToughSession(hub *Hub, client *remoteClient, clean bool) *ToughSession {
	sess := ToughSession{
		hub:           hub,
		clean:         clean,
		outChan:       make(chan *packet.Message),
		holdPackets:   make(map[uint16]packet.Packet),
		currentClient: client,
	}
	sender := NewSender(&sess)
	go sender.Run()
	return &sess
}

func (ts *ToughSession) GetOutChan() chan *packet.Message {
	return ts.outChan
}

func (ts *ToughSession) PileupMsg(msg *packet.Message) error {
	ts.outChan <- msg
	return nil
}

func (ts *ToughSession) HoldPacket(pkt packet.Packet) error {
	//缓存10个包，多余保存到数据库,或者考虑服务器问题，缓存数据怎么办？
	id, ok := packet.PacketID(pkt)
	if ok {
		ts.holdPackets[id] = pkt
	}
	return nil
}

func (ts *ToughSession) LookupPacket(id uint16) (packet.Packet, error) {
	return ts.holdPackets[id], nil
}

func (ts *ToughSession) PublishMsg(pkt *packet.PublishPacket) error {

	ts.currentClient.send(pkt)
	return nil
}

func (ts *ToughSession) GetSubscription() string {
	return ""
}

func (ts ToughSession) AddSubscription(pkt *packet.Subscription) error {
	store := ts.hub.getStore()
	c := store.DB("sunqtt").C("subscriptions")
	err := c.Insert(pkt)
	if err != nil {

	}

	return nil
}

func (ts ToughSession) Save(oid string) error {
	store := ts.hub.getStore()
	c := store.DB("sunqtt").C("sessions")
	ts.oid = oid
	err := c.Insert(&ts)
	if err != nil {

	}

	return nil
}
