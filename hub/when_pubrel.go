package hub

import (
	"fmt"

	"github.com/gomqtt/packet"
)

// handle an incoming PubrelPacket
func (rc *remoteClient) whenPubrel(packetID uint16) error {
	fmt.Println("pubrel...")
	// get packet from store
	pkt, err := rc.session.LookupPacket(packetID)
	if err != nil {
		// return rc.die(err, true)
	}

	// get packet from store
	publish, ok := pkt.(*packet.PublishPacket)
	if !ok {
		return nil // ignore a wrongly sent PubrelPacket
	}

	pubcomp := packet.NewPubcompPacket()
	pubcomp.PacketID = publish.PacketID
	fmt.Println("publish complete...")
	// acknowledge PublishPacket
	err = rc.send(pubcomp)
	if err != nil {
		// return rc.die(err, false)
	}

	// remove packet from store
	// err = c.session.DeletePacket(packetID)
	// if err != nil {
	// 	return c.die(err, true)
	// }

	// publish packet to others
	err = rc.hub.Publish(&publish.Message)
	if err != nil {
		// return rc.die(err, true)
	}

	return nil
}
