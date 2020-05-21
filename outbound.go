package freeswitchesl

import (
	"context"
	"errors"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"log"
	"net"
	"time"
)

type OutboundHandler func(conn *Conn, connectResponse *RawResponse)

/*
 * TODO: Review if we should have a rate limiting facility to prevent DoS attacks
 * For our use it should be fine since we only want to listen on localhost
 */
func ListenAndServe(address string, handler OutboundHandler) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Printf("Listenting for new ESL connections on %s\n", listener.Addr().String())
	for {
		c, err := listener.Accept()
		if err != nil {
			break
		}

		log.Printf("New outbound connection from %s\n", c.RemoteAddr().String())
		conn := newConnection(c)
		// Does not call the handler directly to ensure closing cleanly
		go conn.outboundHandle(handler)
	}
	log.Println("Outbound server shutting down")
	return errors.New("connection closed")
}

func (c *Conn) outboundHandle(handler OutboundHandler) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	response, err := c.sendCommand(ctx, command.Connect{})
	cancel()
	if err != nil {
		log.Printf("Error connecting to %s error %s", c.conn.RemoteAddr().String(), err.Error())
		// Try closing cleanly first
		c.Close()
		return
	}
	handler(c, response)
	c.Close()
}