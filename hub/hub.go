package hub

import (
	"fmt"
	"net"

	"github.com/gomqtt/packet"
	"github.com/showntop/sun-broker/store"
	"gopkg.in/mgo.v2"
)

type Hub struct {
	Sessions map[string]Session //使用用户uuid作为session主键
	Store    *mgo.Database
}

var hub Hub

func init() {
	///初始化存储装置
	store := store.NewStore(1)
	hub = Hub{
		sessions: make(map[string]Session),
		store:    store,
	}
}

func Mount(conn net.Conn) {
	//验证授权用户
	//把client插入hub
	NewRemoteClient(&hub, conn)
	// hub.Regist(client)
}

func (h *Hub) Seed(sess Session, sessid string) error {
	h.sessions[sessid] = sess
	return nil
}

func (h *Hub) Distribute(msg *packet.Message) error {
	// store := h.store
	// c := store.DB("sunqtt").C("subscriptions")
	// subs, err := c.find()
	// if err != nil {

	// }

	for _, session := range h.sessions {
		// fmt.Println(k, v)
		fmt.Println(session)
		session.PileupMsg(msg)
	}

	return nil
}
