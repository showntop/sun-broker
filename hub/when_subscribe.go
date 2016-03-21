package hub

import (
	"fmt"

	"github.com/gomqtt/packet"
)

// handle an incoming SubscribePacket
func (rc *remoteClient) whenSubscribe(pkt *packet.SubscribePacket) error {
	suback := packet.NewSubackPacket()
	suback.ReturnCodes = make([]byte, len(pkt.Subscriptions))
	suback.PacketID = pkt.PacketID

	// var retainedMessages []*packet.Message
	fmt.Println("sub...")
	fmt.Println(pkt.Subscriptions)
	for i, subscription := range pkt.Subscriptions {
		// save subscription in session
		err := rc.session.AddSubscription(&subscription)
		if err != nil {
			// return c.die(err, true)
		}

		// subscribe client to queue
		// msgs, err := rc.broker.Backend.Subscribe(c, subscription.Topic)
		// if err != nil {
		//  return c.die(err, true)
		// }

		// cache retained messages
		// retainedMessages = append(retainedMessages, msgs...)

		// save granted qos
		suback.ReturnCodes[i] = subscription.QOS
	}

	// send suback
	err := rc.send(suback)
	if err != nil {
		// return c.die(err, false)
	}

	// send messages
	// for _, msg := range retainedMessages {
	//  rc.out <- msg
	// }

	return nil
}
