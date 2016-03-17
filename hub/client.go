package hub

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/gomqtt/packet"
)

const (
	clientConnecting byte = iota
	clientConnected
	clientDisconnected
)

type remoteClient struct {
	hub     *Hub
	conn    net.Conn
	session Session
	out     chan *packet.Message
}

// newRemoteClient takes over a connection and returns a remoteClient
func NewRemoteClient(hub *Hub, conn net.Conn) *remoteClient {
	c := &remoteClient{
		hub:  hub,
		conn: conn,
		out:  make(chan *packet.Message),
	}

	go c.loopProc()

	return c
}

func (rc *remoteClient) loopProc() {

	// set initial read timeout
	// rc.conn.SetReadTimeout('')

	for {
		// get next packet from connection
		pkt, err := rc.Receive()
		if err != nil {
			// if rc.state.get() == clientDisconnected {
			// 	return c.die(nil, false)
			// }

			// // die on any other error
			// return c.die(err, false)
		}
		fmt.Println(pkt)
		// c.log("%s - Received: %s", c.Context().Get("uuid"), pkt.String())

		switch pkt.(type) {
		case *packet.ConnectPacket:
			connect, ok := pkt.(*packet.ConnectPacket)
			if !ok {
				// return c.die(fmt.Errorf("expected connect"), true)
			}
			err = rc.whenConnect(connect)
		case *packet.SubscribePacket:
			subcription, ok := pkt.(*packet.SubscribePacket)
			if !ok {
				//停掉客户端
				// return rc.die(fmt.Errorf("expected connect"), true)
			}

			err = rc.whenSubscribe(subcription)
		// case *packet.UnsubscribePacket:
		// 	err = rc.whenUnsubscribe(pkt)
		case *packet.PublishPacket:
			publish, ok := pkt.(*packet.PublishPacket)
			if !ok {
				//停掉客户端
				// return rc.die(fmt.Errorf("expected connect"), true)
			}
			err = rc.whenPublish(publish)
		// case *packet.PubackPacket:
		// 	err = rc.whenPubackAndPubcomp(pkt.PacketID)
		// case *packet.PubcompPacket:
		// 	err = rc.whenPubackAndPubcomp(pkt.PacketID)
		// case *packet.PubrecPacket:
		// 	err = rc.whenPubrec(pkt.PacketID)
		// case *packet.PubrelPacket:
		// 	err = rc.whenPubrel(pkt.PacketID)
		case *packet.PingreqPacket:
			err = rc.whenPingreq()
			// case *packet.DisconnectPacket:
			// 	err = rc.whenDisconnect()
		}

		// return eventual error
		// if err != nil {
		// 	return err // error has already been cleaned
		// }
	}
}

func (rc *remoteClient) whenPingreq() error {
	err := rc.send(packet.NewPingrespPacket())
	if err != nil {
		// return c.die(err, false)
	}

	return nil
}

// sends packet
func (rc *remoteClient) send(pkt packet.Packet) error {

	// rc.sMutex.Lock()
	// defer c.sMutex.Unlock()

	// allocate buffer
	buf := make([]byte, pkt.Len())

	// encode packet
	_, err := pkt.Encode(buf)
	if err != nil {
		rc.conn.Close()
		return errors.New("text") //newTransportError(EncodeError, err)
	}

	// write buffer to connection
	bytesWritten, err := rc.conn.Write(buf)
	if err != nil {
		rc.conn.Close()
		// return newTransportError(NetworkError, err)
	}
	fmt.Println(bytesWritten)
	// increment write counter
	// atomic.AddInt64(&c.writeCounter, int64(bytesWritten))

	// return nil
	return nil
}

func (rc *remoteClient) whenConnect(pkt *packet.ConnectPacket) error {
	connack := packet.NewConnackPacket()
	connack.ReturnCode = packet.ConnectionAccepted
	connack.SessionPresent = false

	// authenticate

	// check authentication

	// set state
	// c.state.set(clientConnected)

	// set keep alive
	// if pkt.KeepAlive > 0 {
	// 	c.conn.SetReadTimeout(time.Duration(pkt.KeepAlive) * 1500 * time.Millisecond)
	// } else {
	// 	c.conn.SetReadTimeout(0)
	// }

	// retrieve session
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
	sender := sender{rc}
	go sender.Run()

	// // retrieve stored packets
	// packets, err := c.session.AllPackets(outgoing)
	// if err != nil {
	// 	return c.die(err, true)
	// }

	// // resend stored packets
	// for _, pkt := range packets {
	// 	publish, ok := pkt.(*packet.PublishPacket)
	// 	if ok {
	// 		// set the dup flag on a publish packet
	// 		publish.Dup = true
	// 	}

	// 	err = c.send(pkt)
	// 	if err != nil {
	// 		return c.die(err, false)
	// 	}
	// }

	// // get stored subscriptions
	// subs, err := c.session.AllSubscriptions()
	// if err != nil {
	// 	return c.die(err, true)
	// }

	// // restore subscriptions
	// for _, sub := range subs {
	// 	// TODO: Handle incoming retained messages.
	// 	c.broker.Backend.Subscribe(c, sub.Topic)
	// }

	return nil
}

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
		// 	return c.die(err, true)
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
	// 	rc.out <- msg
	// }

	return nil
}

// Receive will read from the underlying connection and return a fully read
// packet. It will return an Error if there was an error while decoding or
// reading from the underlying connection.
//
// Note: Only one goroutine can Receive at the same time.
func (rc *remoteClient) Receive() (packet.Packet, error) {
	// c.rMutex.Lock()
	// defer c.rMutex.Unlock()

	// initial detection length
	detectionLength := 2

	for {
		// check length
		if detectionLength > 5 {
			rc.conn.Close()
			return nil, errors.New("fjdsalkfjlksf") //errors.New("fjdsalkfjlksf")//newTransportError(DetectionError, ErrDetectionOverflow)
		}

		// try read detection bytes
		reader := bufio.NewReader(rc.conn)

		header, err := reader.Peek(detectionLength)
		if err == io.EOF && len(header) == 0 {
			// only if Peek returned no bytes the close is expected
			rc.conn.Close()
			return nil, errors.New("fjdsalkfjlksf") //errors.New("fjdsalkfjlksf")//newTransportError(ConnectionClose, err)
		} else if opError, ok := err.(*net.OpError); ok && opError.Timeout() {
			// the read timed out
			rc.conn.Close()
			return nil, errors.New("fjdsalkfjlksf") //newTransportError(NetworkError, ErrReadTimeout)
		} else if err != nil {
			rc.conn.Close()
			return nil, errors.New("fjdsalkfjlksf") //newTransportError(NetworkError, err)
		}

		// detect packet
		packetLength, packetType := packet.DetectPacket(header)

		// on zero packet length:
		// increment detection length and try again
		if packetLength <= 0 {
			detectionLength++
			continue
		}

		// check read limit
		// if c.readLimit > 0 && int64(packetLength) > c.readLimit {
		// 	c.conn.Close()
		// 	return nil, errors.New("fjdsalkfjlksf") //newTransportError(NetworkError, ErrReadLimitExceeded)
		// }

		// create packet
		pkt, err := packetType.New()
		if err != nil {
			rc.conn.Close()
			return nil, errors.New("fjdsalkfjlksf") //newTransportError(DetectionError, err)
		}

		// allocate buffer
		buf := make([]byte, packetLength)

		// read whole packet
		bytesRead, err := io.ReadFull(reader, buf)
		if opError, ok := err.(*net.OpError); ok && opError.Timeout() {
			// the read timed out
			rc.conn.Close()
			return nil, errors.New("fjdsalkfjlksf") //newTransportError(NetworkError, ErrReadTimeout)
		} else if err != nil {
			rc.conn.Close()

			// even if EOF is returned we consider it an network error
			return nil, errors.New("fjdsalkfjlksf") //newTransportError(NetworkError, err)
		}
		fmt.Println(bytesRead)
		// decode buffer
		_, err = pkt.Decode(buf)
		if err != nil {
			rc.conn.Close()
			return nil, errors.New("fjdsalkfjlksf") //newTransportError(DecodeError, err)
		}

		// increment counter
		// atomic.AddInt64(&c.readCounter, int64(bytesRead))

		// // reset timeout
		// c.resetTimeout()

		return pkt, nil
	}
}
