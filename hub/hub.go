package hub

import (
	"net"

	"github.com/showntop/sun-broker/store"
	"gopkg.in/mgo.v2"
)

type Hub struct {
	sessions map[string]*Session
	store    *mgo.Session
}

var hub Hub

func init() {
	///初始化存储装置
	store := store.NewStore(1)
	hub = Hub{
		sessions: make(map[string]*Session),
		store:    store,
	}
}

func (h *Hub) getStore() *mgo.Session {
	return h.store
}

func Mount(conn net.Conn) {
	//验证授权用户
	//把client插入hub
	NewRemoteClient(&hub, conn)
	// hub.Regist(client)
}

func (h *Hub) Seed(sess Session, sessid string) error {
	h.sessions[sessid] = &sess
	return nil
}
