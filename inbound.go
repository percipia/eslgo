package freeswitchesl

import (
	"context"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"log"
	"net"
)

func Dial(address, password string, onDisconnect func()) (*Conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	connection := newConnection(c, false)

	// Inbound only handlers
	go connection.authLoop(command.Auth{Password: password})
	go connection.disconnectLoop(onDisconnect)

	return connection, nil
}

func (c *Conn) disconnectLoop(onDisconnect func()) {
	select {
	case <-c.responseChannels[TypeDisconnect]:
		c.Close()
		defer onDisconnect()
		return
	case <-c.runningContext.Done():
		return
	}
}

func (c *Conn) authLoop(auth command.Auth) {
	for {
		select {
		case <-c.responseChannels[TypeAuthRequest]:
			response, err := c.SendCommand(context.Background(), auth)
			if err != nil {
				log.Printf("Failed to auth %e\n", err)
				return
			}
			if !response.IsOk() {
				// Try to gracefully disconnect
				log.Printf("Failed to auth %#v\n", response)
				_, _ = c.SendCommand(context.Background(), command.Exit{})
				return
			} else {
				log.Printf("Sucessfully authenticated %s\n", c.conn.RemoteAddr())
			}
		case <-c.runningContext.Done():
			return
		}
	}
}
