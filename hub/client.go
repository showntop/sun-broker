package hub

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/gomqtt/packet"
)

//客户端当前连接状态
const (
	clientConnecting byte = iota
	clientConnected
	clientDisconnected
)

type remoteClient struct {
	rMutex  sync.Mutex
	sMutex  sync.Mutex
	hub     *Hub
	conn    net.Conn
	session Session
	state   byte
}

// newRemoteClient takes over a connection and returns a remoteClient
func NewRemoteClient(hub *Hub, conn net.Conn) *remoteClient {
	c := &remoteClient{
		hub:  hub,
		conn: conn,
	}

	go c.run()

	return c
}

func (rc *remoteClient) run() {

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
		if pkt == nil {
			panic(pkt)
		}
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
		case *packet.PubrecPacket:
			pubrecPacket, ok := pkt.(*packet.PubrecPacket)
			if !ok {
				//停掉客户端
				// return rc.die(fmt.Errorf("expected connect"), true)
			}

			err = rc.whenPubrec(pubrecPacket.PacketID)
		case *packet.PubrelPacket:
			pubrelPacket, ok := pkt.(*packet.PubrelPacket)
			if !ok {
				//停掉客户端
				// return rc.die(fmt.Errorf("expected connect"), true)
			}

			err = rc.whenPubrel(pubrelPacket.PacketID)
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

// Receive will read from the underlying connection and return a fully read
// packet. It will return an Error if there was an error while decoding or
// reading from the underlying connection.
//
// Note: Only one goroutine can Receive at the same time.
func (rc *remoteClient) Receive() (packet.Packet, error) {
	rc.rMutex.Lock()
	defer rc.rMutex.Unlock()

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

// sends packet
func (rc *remoteClient) send(pkt packet.Packet) error {
	rc.sMutex.Lock()
	defer rc.sMutex.Unlock()

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
