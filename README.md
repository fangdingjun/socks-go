socks-go
=======

A socks server implement by golang, support socks4/4a, socks5.


usage
====
Usage example:

    import "github.com/fangdingjun/socks"

    fucn main(){
        l, _ := net.Listen("tcp", ":1080")
        for {
            conn, _ := l.Accept()
            s := socks.SocksConn{conn}
            go s.Serve()
        }
    }
