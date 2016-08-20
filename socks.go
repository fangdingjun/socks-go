package socks

import (
	"io"
	"log"
	"net"
	"time"
)

const (
	socks4Version  = 0x04
	socks5Version  = 0x05
	cmdConnect     = 0x01
	addrTypeIPv4   = 0x01
	addrTypeDomain = 0x03
	addrTypeIPv6   = 0x04
)

type dialFunc func(network, addr string) (net.Conn, error)

// SocksConn present a client connection
type SocksConn struct {
	ClientConn net.Conn
	Dial       dialFunc
}

// Serve serve the client
func (s *SocksConn) Serve() {
	buf := make([]byte, 1)

	// read version
	io.ReadFull(s.ClientConn, buf)

	dial := s.Dial
	if s.Dial == nil {
		d := net.Dialer{Timeout: 10 * time.Second}
		dial = d.Dial
	}

	switch buf[0] {
	case socks4Version:
		s4 := socks4Conn{clientConn: s.ClientConn, dial: dial}
		s4.Serve()
	case socks5Version:
		s5 := socks5Conn{clientConn: s.ClientConn, dial: dial}
		s5.Serve()
	default:
		log.Printf("error version %s", buf[0])
		s.ClientConn.Close()
	}
}
