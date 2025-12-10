/*
 * Copyright (c) 2020 Percipia
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Contributor(s):
 * Andrew Querol <aquerol@percipia.com>
 */
package eslgo

import (
	"context"
	"errors"
	"github.com/percipia/eslgo/command"
	"net"
	"time"
)

type OutboundHandler func(ctx context.Context, conn *Conn, connectResponse *RawResponse)

// OutboundOptions - Used to open a new listener for outbound ESL connections from FreeSWITCH
type OutboundOptions struct {
	Options                       // Generic common options to both Inbound and Outbound Conn
	Network         string        // The network type to listen on, should be tcp, tcp4, or tcp6
	ConnectTimeout  time.Duration // How long should we wait for FreeSWITCH to respond to our "connect" command. 5 seconds is a sane default.
	ConnectionDelay time.Duration // How long should we wait after connection to start sending commands. 25ms is the recommended default otherwise we can close the connection before FreeSWITCH finishes starting it on their end. https://github.com/signalwire/freeswitch/pull/636
}

// DefaultOutboundOptions - The default options used for creating the outbound connection
var DefaultOutboundOptions = OutboundOptions{
	Options:         DefaultOptions,
	Network:         "tcp",
	ConnectTimeout:  5 * time.Second,
	ConnectionDelay: 25 * time.Millisecond,
}

/*
 * TODO: Review if we should have a rate limiting facility to prevent DoS attacks
 * For our use it should be fine since we only want to listen on localhost
 */
// ListenAndServe - Open a new listener for outbound ESL connections from FreeSWITCH on the specified address with the provided connection handler
func ListenAndServe(address string, handler OutboundHandler) error {
	return DefaultOutboundOptions.ListenAndServe(address, handler)
}

// ListenAndServe - Open a new listener for outbound ESL connections from FreeSWITCH with provided options and handle them with the specified handler
func (opts OutboundOptions) ListenAndServe(address string, handler OutboundHandler) error {
	listener, err := net.Listen(opts.Network, address)
	if err != nil {
		return err
	}
	if opts.Logger != nil {
		opts.Logger.Info("Listening for new ESL connections on %s\n", listener.Addr().String())
	}
	for {
		c, err := listener.Accept()
		if err != nil {
			break
		}
		conn := newConnection(c, true, opts.Options)

		conn.logger.Info("New outbound connection from %s\n", c.RemoteAddr().String())
		go conn.dummyLoop()
		// Does not call the handler directly to ensure closing cleanly
		go conn.outboundHandle(handler, opts.ConnectionDelay, opts.ConnectTimeout)
	}

	if opts.Logger != nil {
		opts.Logger.Info("Outbound server shutting down")
	}
	return errors.New("connection closed")
}

func (c *Conn) outboundHandle(handler OutboundHandler, connectionDelay, connectTimeout time.Duration) {
	ctx, cancel := context.WithTimeout(c.runningContext, connectTimeout)
	response, err := c.SendCommand(ctx, command.Connect{})
	cancel()
	if err != nil {
		c.logger.Warn("Error connecting to %s error %s", c.conn.RemoteAddr().String(), err.Error())
		// Try closing cleanly first
		c.Close() // Not ExitAndClose since this error connection is most likely from communication failure
		return
	}
	handler(c.runningContext, c, response)
	// XXX This is ugly, the issue with short lived async sockets on our end is if they complete too fast we can actually
	// close the connection before FreeSWITCH is in a state to close the connection on their end. 25ms is an magic value
	// found by testing to have no failures on my test system. I started at 1 second and reduced as far as I could go.
	// TODO This actually may be fixed: https://github.com/signalwire/freeswitch/pull/636
	time.Sleep(connectionDelay)
	c.ExitAndClose()
}

func (c *Conn) dummyLoop() {
	select {
	case <-c.responseChannels[TypeDisconnect]:
		c.logger.Info("Disconnect outbound connection", c.conn.RemoteAddr())
		if c.closeDelay >= 0 {
			time.AfterFunc(c.closeDelay*time.Second, func() {
				c.Close()
			})
		}
	case <-c.responseChannels[TypeAuthRequest]:
		c.logger.Debug("Ignoring auth request on outbound connection", c.conn.RemoteAddr())
	case <-c.runningContext.Done():
		return
	}
}
