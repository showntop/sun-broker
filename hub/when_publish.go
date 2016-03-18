package hub

import (
	"fmt"

	"github.com/gomqtt/packet"
)

func (rc *remoteClient) whenPublish(publish *packet.PublishPacket) error {

	fmt.Println("publish...")
	if publish.Message.QOS == 1 {
		fmt.Println("qos1...")

		puback := packet.NewPubackPacket()
		puback.PacketID = publish.PacketID

		// acknowledge qos 1 publish
		err := rc.send(puback)
		if err != nil {
			// return rc.die(err, false)
		}
	}

	if publish.Message.QOS == 2 {
		fmt.Println("qos2...")

		// store packet
		err := rc.session.SaveOutPacket(incoming, publish)
		if err != nil {
			// return rc.die(err, true)
		}

		pubrec := packet.NewPubrecPacket()
		pubrec.PacketID = publish.PacketID

		// signal qos 2 publish
		err = rc.send(pubrec)
		if err != nil {
			// return c.die(err, false)
		}
	}

	if publish.Message.QOS <= 1 {
		// publish packet to others
		// err := rc.broker.Backend.Publish(c, &publish.Message)
		// if err != nil {
		// 	// return rc.die(err, true)
		// }
	}

	return nil
}
