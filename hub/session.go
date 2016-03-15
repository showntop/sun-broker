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
