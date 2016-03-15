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

func NewStore(driverType int) *mgo.Session {
	return ConnectMongo("mongodb://127.0.0.1")
}
