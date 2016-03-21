package store

import (
	"gopkg.in/mgo.v2"
)

type subscription struct {
	Id    string
	topic string
}

type SessionStore struct {
	Oid           string
	subscriptions []subscription
}

func NewStore(driverType int) *mgo.Database {
	///从配置文件中
	session := ConnectMongo("mongodb://127.0.0.1")
	return session.DB("sunqtt")
}
