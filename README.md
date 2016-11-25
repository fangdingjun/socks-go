socks-go
=======

A socks server implemented by golang, support socks4/4a, socks5.

Only support connect command now.

usage
====
Usage example:

    import socks "github.com/fangdingjun/socks-go"

    func main(){
        l, _ := net.Listen("tcp", ":1080")
        for {
            conn, _ := l.Accept()
            s := socks.SocksConn{ClientConn: conn, Dial: nil} // Dial is a function which dial to the upstream server
            go s.Serve() // serve the socks request
        }
    }
