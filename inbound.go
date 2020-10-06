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
package freeswitchesl

import (
	"context"
	"fmt"
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

	// First auth
	<-connection.responseChannels[TypeAuthRequest]
	err = connection.doAuth(connection.runningContext, command.Auth{Password: password})
	if err != nil {
		// Try to gracefully disconnect
		log.Printf("Failed to auth %e\n", err)
		_, _ = connection.SendCommand(connection.runningContext, command.Exit{})
	} else {
		log.Printf("Sucessfully authenticated %s\n", connection.conn.RemoteAddr())
	}

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
			err := c.doAuth(c.runningContext, auth)
			if err != nil {
				// Try to gracefully disconnect
				log.Printf("Failed to auth %e\n", err)
				_, _ = c.SendCommand(c.runningContext, command.Exit{})
			} else {
				log.Printf("Sucessfully authenticated %s\n", c.conn.RemoteAddr())
			}
		case <-c.runningContext.Done():
			return
		}
	}
}

func (c *Conn) doAuth(ctx context.Context, auth command.Auth) error {
	response, err := c.SendCommand(ctx, auth)
	if err != nil {
		return err
	}
	if !response.IsOk() {
		return fmt.Errorf("failed to auth %#v", response)
	}
	return nil
}
