package hub

import (
	// "fmt"

	"github.com/gomqtt/packet"
)

func (rc *remoteClient) whenConnect(pkt *packet.ConnectPacket) error {
	connack := packet.NewConnackPacket()
	connack.ReturnCode = packet.ConnectionAccepted
	connack.SessionPresent = false

	// authenticate

	// check authentication

	//set state
	rc.state.set(clientConnected)

	// set keep alive
	// if pkt.KeepAlive > 0 {
	//  rc.conn.SetReadTimeout(time.Duration(pkt.KeepAlive) * 1500 * time.Millisecond)
	// } else {
	//  rc.conn.SetReadTimeout(0)
	// }

	// retrieve session
	//from the config level
	sess := NewToughSession(rc.hub, rc, pkt.CleanSession)
	err := rc.hub.Seed(sess, pkt.ClientID)
	if err != nil {
		// return rc.die(err, true)
	}
	sess.Save(pkt.ClientID)
	// set session present
	connack.SessionPresent = !pkt.CleanSession

	// assign session
	rc.session = sess

	// save will if present
	if pkt.Will != nil {
		err = nil //rc.session.SaveWill(pkt.Will)
		if err != nil {
			// return c.die(err, true)
		}
	}

	// send connack
	err = rc.send(connack)
	if err != nil {
		// return rc.die(err, false)
	}

	// start sender
	// go c.sender
	// sender := sender{rc}
	// go sender.Run()

	// // retrieve stored packets
	// packets, err := c.session.AllPackets(outgoing)
	// if err != nil {
	//  return c.die(err, true)
	// }

	// // resend stored packets
	// for _, pkt := range packets {
	//  publish, ok := pkt.(*packet.PublishPacket)
	//  if ok {
	//      // set the dup flag on a publish packet
	//      publish.Dup = true
	//  }

	//  err = c.send(pkt)
	//  if err != nil {
	//      return c.die(err, false)
	//  }
	// }

	// // get stored subscriptions
	// subs, err := c.session.AllSubscriptions()
	// if err != nil {
	//  return c.die(err, true)
	// }

	// // restore subscriptions
	// for _, sub := range subs {
	//  // TODO: Handle incoming retained messages.
	//  c.broker.Backend.Subscribe(c, sub.Topic)
	// }

	return nil
}
