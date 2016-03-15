package store

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

func ConnectMongo(url string) *mgo.Session {
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println(err)
	}
	return session
}
