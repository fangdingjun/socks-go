package main

import (
	"github.com/fangdingjun/socks"
	"log"
	"net"
)

func main() {
	conn, err := net.Listen("tcp", ":1080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := conn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("connected from %s", c.RemoteAddr())
		s := socks.SocksConn{ClientConn: c}
		go s.Serve()
	}
}
