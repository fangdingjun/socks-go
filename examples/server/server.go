package main

import (
	"log"
	"net"
	"time"

	socks "github.com/fangdingjun/socks-go"
)

func main() {
	conn, err := net.Listen("tcp", ":1080")
	if err != nil {
		log.Fatal(err)
	}
	a := socks.PasswordAuth{Username: "admin", Password: "passwd"}
	for {
		c, err := conn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("connected from %s", c.RemoteAddr())

		d := net.Dialer{Timeout: 10 * time.Second}
		s := socks.Conn{Conn: c, Dial: d.Dial, Auth: &a}
		go s.Serve()
	}
}
