package hub

import (
	// "github.com/showntop/sun-broker/store"
	"github.com/gomqtt/packet"
)

type Session interface {
	AddSubscription(*packet.Subscription) error
}

type ToughSession struct {
	oid           string
	hub           *Hub
	subscriptions []string
	offlineMsgs   []string
	outPackets    map[string]string
	inPackets     map[string]string

	will *packet.Message

	currentClient remoteClient
	clean         bool
}

func NewToughSession(hub *Hub, client *remoteClient, clean bool) ToughSession {
	return ToughSession{
		hub:   hub,
		clean: clean,
	}
}

func (ts *Session) CurrentClient() {
	return ts.currentClient
}

func (ts *ToughSession) LookupPacket(id uint16) error {

}

func (ts *ToughSession) SavePacket(pkt packet.Packet) error {
	//缓存10个包，多余保存到数据库,或者考虑服务器问题，缓存数据怎么办？
	id, ok := packet.PacketID(pkt)
	if ok {
		ts.outPackets[id] = pkt
	}

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
