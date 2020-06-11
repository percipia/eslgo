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
		conn := newConnection(c, true)
		go conn.dummyLoop()
		// Does not call the handler directly to ensure closing cleanly
		go conn.outboundHandle(handler)
	}
	log.Println("Outbound server shutting down")
	return errors.New("connection closed")
}

func (c *Conn) outboundHandle(handler OutboundHandler) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	response, err := c.SendCommand(ctx, command.Connect{})
	cancel()
	if err != nil {
		log.Printf("Error connecting to %s error %s", c.conn.RemoteAddr().String(), err.Error())
		// Try closing cleanly first
		c.Close()
		return
	}
	handler(c, response)
	// XXX This is ugly, the issue with short lived async sockets on our end is if they complete too fast we can actually
	// close the connection before FreeSWITCH is in a state to close the connection on their end. 25ms is an magic value
	// found by testing to have no failures on my test system. I started at 1 second and reduced as far as I could go.
	// TODO We should open a bug report on the FreeSWITCH GitHub at some point and remove this when fixed.
	// TODO This actually may be fixed: https://github.com/signalwire/freeswitch/pull/636
	time.Sleep(25 * time.Millisecond)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	_, _ = c.SendCommand(ctx, command.Exit{})
	cancel()
	c.Close()
}

func (c *Conn) dummyLoop() {
	select {
	case <-c.responseChannels[TypeDisconnect]:
		log.Println("Disconnect outbound connection", c.conn.RemoteAddr())
		c.Close()
	case <-c.responseChannels[TypeAuthRequest]:
		log.Println("Ignoring auth request on outbound connection", c.conn.RemoteAddr())
	case <-c.runningContext.Done():
		return
	}
}
