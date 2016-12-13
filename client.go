package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// Client is a net.Conn with socks5 support
type Client struct {
	net.Conn
	// socks5 username
	Username string
	// socks5 password
	Password      string
	handshakeDone bool
	connected     bool
}

func (sc *Client) handShake() error {
	if sc.handshakeDone {
		return nil
	}

	// password auth or none
	if _, err := sc.Conn.Write([]byte{socks5Version, 0x02, 0x00, 0x02}); err != nil {
		return err
	}

	buf := make([]byte, 2)
	if _, err := io.ReadFull(sc.Conn, buf); err != nil {
		return err
	}

	if buf[0] != socks5Version {
		return fmt.Errorf("error socks version %s", buf[0])
	}

	if buf[1] != 0x00 && buf[1] != 0x02 {
		return fmt.Errorf("server return with code %s", buf[1])
	}

	if buf[1] == 0x00 {
		sc.handshakeDone = true
		return nil
	}

	// password auth

	buf = make([]byte, 3+len(sc.Username)+len(sc.Password))

	buf[0] = 0x01                                       // auth protocol version
	buf[1] = byte(len(sc.Username))                     // username length
	copy(buf[2:], []byte(sc.Username))                  // username
	buf[2+len(sc.Username)] = byte(len(sc.Password))    // password length
	copy(buf[3+len(sc.Username):], []byte(sc.Password)) //password

	if _, err := sc.Conn.Write(buf); err != nil {
		return err
	}

	if _, err := sc.Conn.Read(buf[:2]); err != nil {
		return err
	}

	if buf[0] != 0x01 {
		return fmt.Errorf("unexpected auth protocol version %v", buf[0])
	}

	// password auth success
	if buf[1] == 0x00 {
		sc.handshakeDone = true
		return nil
	}

	return fmt.Errorf("password rejected")
}

// Connect connects to socks5 server and handshake
func (sc *Client) Connect(host string, port uint16) error {
	if !sc.handshakeDone {
		if err := sc.handShake(); err != nil {
			return err
		}
	}

	if sc.connected {
		return nil
	}

	l := 4 + len(host) + 1 + 2
	buf := make([]byte, l)
	buf[0] = socks5Version
	buf[1] = cmdConnect
	buf[2] = 0x00
	buf[3] = addrTypeDomain
	buf[4] = byte(len(host))

	copy(buf[5:5+len(host)], []byte(host))

	binary.BigEndian.PutUint16(buf[l-2:l], port)

	if _, err := sc.Conn.Write(buf); err != nil {
		return err
	}

	buf1 := make([]byte, 128)

	if _, err := io.ReadAtLeast(sc.Conn, buf1, 10); err != nil {
		return err
	}

	if buf1[0] != socks5Version {
		return fmt.Errorf("error socks version %d", buf1[0])
	}

	if buf1[1] != 0x00 {
		return fmt.Errorf("server error code %d", buf1[1])
	}

	sc.connected = true
	return nil
}

// Read read from the underlying connection
func (sc *Client) Read(b []byte) (int, error) {
	if !sc.connected {
		return 0, fmt.Errorf("call connect first")
	}
	return sc.Conn.Read(b)
}

// Write write data to underlying connection
func (sc *Client) Write(b []byte) (int, error) {
	if !sc.connected {
		return 0, fmt.Errorf("call connect first")
	}
	return sc.Conn.Write(b)
}

// Close close the underlying connection
func (sc *Client) Close() error {
	return sc.Conn.Close()
}
