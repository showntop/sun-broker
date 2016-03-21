package hub

import (
	// "fmt"

	"github.com/gomqtt/packet"
)

func (rc *remoteClient) whenPingreq() error {
	err := rc.send(packet.NewPingrespPacket())
	if err != nil {
		// return rc.die(err, false)
	}

	return nil
}
