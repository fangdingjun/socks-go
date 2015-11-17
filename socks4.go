package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type socks4Conn struct {
	server_conn net.Conn
	client_conn net.Conn
}

func (s4 *socks4Conn) Serve() {
	defer s4.Close()
	if err := s4.processRequest(); err != nil {
		log.Println(err)
		return
	}
}

func (s4 *socks4Conn) Close() {
	if s4.client_conn != nil {
		s4.client_conn.Close()
	}
	if s4.server_conn != nil {
		s4.server_conn.Close()
	}
}

func (s4 *socks4Conn) forward() {
	go func() {
		io.Copy(s4.client_conn, s4.server_conn)
	}()

	io.Copy(s4.server_conn, s4.client_conn)
}

func (s4 *socks4Conn) processRequest() error {
	// version has already read out by socksConn.Serve()
	// process command and target here

	buf := make([]byte, 128)
	n, err := io.ReadAtLeast(s4.client_conn, buf, 8)
	if err != nil {
		return err
	}

	// only support connect
	if buf[0] != cmdConnect {
		return fmt.Errorf("error command %s", buf[0])
	}

	port := binary.BigEndian.Uint16(buf[1:3])

	ip := net.IP(buf[3:7])

	// NULL-terminated user string
	// jump to NULL character
	var j int
	for j = 7; j < n; j++ {
		if buf[j] == 0x00 {
			break
		}
	}

	host := ip.String()

	// socks4a
	// 0.0.0.x
	if ip[0] == 0x00 && ip[1] == 0x00 && ip[2] == 0x00 && ip[3] != 0x00 {
		j++
		var i = j

		// jump to the end of hostname
		for j = i; j < n; j++ {
			if buf[j] == 0x00 {
				break
			}
		}
		host = string(buf[i:j])
	}

	target := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	// reply user with connect success
	// if dial to target failed, user will receive connection reset
	s4.client_conn.Write([]byte{0x00, 0x5a, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00})

	log.Printf("connecting to %s", target)

	// connect to the target
	s4.server_conn, err = net.Dial("tcp", target)
	if err != nil {
		return err
	}

	// enter data exchange
	s4.forward()

	return nil
}
