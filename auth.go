package socks

import "net"

// AuthService the service to authenticate the user for socks5
type AuthService interface {
	// Authenticate auth the user
	// return true means ok, false means no access
	Authenticate(username, password string, addr net.Addr) bool
}

// default password auth service
type PasswordAuth struct {
	Username string
	Password string
}

func (pa *PasswordAuth) Authenticate(username, password string, addr net.Addr) bool {
	if username == pa.Username && password == pa.Password {
		return true
	}
	return false
}
