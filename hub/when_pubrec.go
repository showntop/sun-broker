package hub

import (
	"fmt"

	"github.com/gomqtt/packet"
)

// handle an incoming PubrelPacket
func (rc *remoteClient) whenPubrec(packetID uint16) error {
	fmt.Println("pubrec...")
	// allocate packet
	pubrel := packet.NewPubrelPacket()
	pubrel.PacketID = packetID

	// overwrite stored PublishPacket with PubrelPacket
	// err := rc.session.SavePacket(outgoing, pubrel)
	// if err != nil {
	// 	return rc.die(err, true)
	// }

	// send packet
	err := rc.send(pubrel)
	if err != nil {
		// return rc.die(err, false)
	}

	return nil
}
