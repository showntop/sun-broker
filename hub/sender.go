package hub

import (
	"github.com/gomqtt/packet"
)

type sender struct {
	session Session
}

func NewSender(sess Session) *sender {
	return &sender{session: sess}
}

func (s *sender) Run() {

	for {
		select {
		// case <-c.tomb.Dying():
		// 	return tomb.ErrDying
		case msg := <-s.session.GetOutChan():
			publish := packet.NewPublishPacket()
			publish.Message = *msg
			publish.PacketID = 10
			// send packet
			err := s.session.PublishMsg(publish)
			if err != nil {
				// return rc.die(err, false)
			}
		}
	}

}
