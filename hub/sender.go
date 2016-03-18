package hub

import (
// "github.com/gomqtt/packet"
)

type sender struct {
	client *remoteClient
}

func (s *sender) Run() {

	// for {
	// 	select {
	// 	// case <-c.tomb.Dying():
	// 	// 	return tomb.ErrDying
	// 	case msg := <-s.client.out:
	// 		publish := packet.NewPublishPacket()
	// 		publish.Message = *msg

	// 		// get stored subscription
	// 		sub, err := s.client.session.LookupSubscription(publish.Message.Topic)
	// 		if err != nil {
	// 			// return rc.die(err, true)
	// 		}

	// 		// check subscription
	// 		if sub == nil {
	// 			// return c.die(fmt.Errorf("subscription not found in session"), true)
	// 		}

	// 		// respect maximum qos
	// 		if publish.Message.QOS > sub.QOS {
	// 			publish.Message.QOS = sub.QOS
	// 		}

	// 		// set packet id
	// 		if publish.Message.QOS > 0 {
	// 			publish.PacketID = c.session.PacketID()
	// 		}

	// 		// store packet if at least qos 1
	// 		if publish.Message.QOS > 0 {
	// 			err := c.session.SavePacket(outgoing, publish)
	// 			if err != nil {
	// 				// return rc.die(err, true)
	// 			}
	// 		}

	// 		// send packet
	// 		err = s.client.send(publish)
	// 		if err != nil {
	// 			// return rc.die(err, false)
	// 		}
	// 	}
	// }

}
