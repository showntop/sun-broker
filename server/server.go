package server

import (
	"net"
	"net/url"
)

// A Server is a local port on which incoming connections can be accepted.
type Server interface {
	// Accept will return the next available connection or block until a
	// connection becomes available, otherwise returns an Error.
	Accept() (net.Conn, error)

	// Close will close the underlying listener and cleanup resources. It will
	// return an Error if the underlying listener didn't close cleanly.
	Close() error

	New(string) error
}

var Srv Server

func init() {

}

// Launch will launch a server based on information extracted from an URL.
func Launch(urlString string) (Server, error) {
	urlParts, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, err //errors.New("url解析错误")
	}
	//使用反射获取相应的服务类型
	switch urlParts.Scheme {
	case "tcp", "mqtt":
		return NewNetServer(urlParts.Host)
	case "tls", "mqtts":
	//	return NewSecureNetServer(urlParts.Host, l.TLSConfig)
	case "ws":
	//	return NewWebSocketServer(urlParts.Host)
	case "wss":
		//	return NewSecureWebSocketServer(urlParts.Host, l.TLSConfig)
	}

	return nil, UnsupportedProtocol //newTransportError(LaunchError, ErrUnsupportedProtocol)
}
