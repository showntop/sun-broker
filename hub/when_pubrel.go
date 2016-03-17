package hub

// handle an incoming PubrelPacket
func (c *remoteClient) processPubrel(packetID uint16) error {
	fmt.Println("pubrel...")
	// get packet from store
	pkt, err := c.session.LookupPacket(incoming, packetID)
	if err != nil {
		return c.die(err, true)
	}

	// get packet from store
	publish, ok := pkt.(*packet.PublishPacket)
	if !ok {
		return nil // ignore a wrongly sent PubrelPacket
	}

	pubcomp := packet.NewPubcompPacket()
	pubcomp.PacketID = publish.PacketID

	// acknowledge PublishPacket
	err = c.send(pubcomp)
	if err != nil {
		return c.die(err, false)
	}

	// remove packet from store
	err = c.session.DeletePacket(incoming, packetID)
	if err != nil {
		return c.die(err, true)
	}

	// publish packet to others
	err = c.broker.Backend.Publish(c, &publish.Message)
	if err != nil {
		return c.die(err, true)
	}

	return nil
}
