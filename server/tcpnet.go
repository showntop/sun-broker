package server

import (
	"net"
)

type NetServer struct {
	listener net.Listener
}

// NewNetServer creates a new TCP server that listens on the provided address.
func NewNetServer(address string) (*NetServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, newServerError(ListentError, err)
	}

	return &NetServer{
		listener: listener,
	}, nil
}

// Accept will return the next available connection or block until a
// connection becomes available, otherwise returns an Error.
func (s *NetServer) Accept() (net.Conn, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, newServerError(ListentError, err)
	}

	return conn, nil
}

// Close will close the underlying listener and cleanup resources. It will
// return an Error if the underlying listener didn't close cleanly.
func (s *NetServer) Close() error {
	err := s.listener.Close()
	if err != nil {
		return newServerError(ListentError, err)
	}

	return nil
}
