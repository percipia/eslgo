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
	"bufio"
	"context"
	"github.com/percipia/eslgo/command"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
)

func TestConn_SendCommand(t *testing.T) {
	server, client := net.Pipe()
	connection := newConnection(client, false, DefaultOptions)
	defer connection.Close()
	defer server.Close()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	serverReader := bufio.NewReader(server)
	defer serverReader.Discard(serverReader.Buffered())

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		response, err := connection.SendCommand(ctx, command.Auth{
			Password: "test1234",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsOk())
		assert.Equal(t, "+OK Job-UUID: c7709e9c-1517-11dc-842a-d3a3942d3d63", response.GetHeader("Reply-Text"))
		wait.Done()
	}()

	// This is sorta lazy, we should be reading until the proper deliminator of \r\n\r\n
	incomingCommand, err := serverReader.ReadString('\r')
	assert.Nil(t, err)
	assert.Equal(t, "auth test1234\r", incomingCommand)

	_, err = server.Write([]byte("Content-Type: command/reply\r\nReply-Text: +OK Job-UUID: c7709e9c-1517-11dc-842a-d3a3942d3d63\r\n\r\n"))
	assert.Nil(t, err)
	wait.Wait()
}

func TestReceiveLoop_UnexpectedEOF(t *testing.T) {
	server, client := net.Pipe()

	disconnected := make(chan struct{})
	opts := DefaultInboundOptions
	opts.Password = "ClueCon"
	opts.OnDisconnect = func() {
		close(disconnected)
	}

	connection := newConnection(client, false, opts.Options)
	go connection.disconnectLoop(opts.OnDisconnect)

	// Send auth/request so eslgo is in a normal state, then send a partial
	// response: headers with Content-Length but a truncated body. When the
	// server closes the connection, io.ReadFull returns "unexpected EOF".
	_, err := server.Write([]byte("Content-Type: auth/request\r\n\r\n"))
	assert.NoError(t, err)

	// Give the connection a moment to process the auth request
	time.Sleep(50 * time.Millisecond)

	// Send a response with Content-Length: 100 but only 10 bytes of body,
	// then close the server side to trigger unexpected EOF in io.ReadFull.
	_, err = server.Write([]byte("Content-Type: command/reply\r\nContent-Length: 100\r\n\r\n0123456789"))
	assert.NoError(t, err)
	server.Close()

	select {
	case <-disconnected:
		// onDisconnect fired — this is the expected behavior
	case <-time.After(5 * time.Second):
		t.Fatal("onDisconnect was not called after unexpected EOF")
	}

	connection.Close()
}

func TestReceiveLoop_NoRaceOnEOF(t *testing.T) {
	server, client := net.Pipe()
	connection := newConnection(client, false, DefaultOptions)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(20 * time.Millisecond)
		server.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(30 * time.Millisecond)
		connection.Close()
	}()

	wg.Wait()
}

// TestReceiveLoop_ReadKeepaliveDeadLink verifies that the read keepalive
// deadline detects a half-open dead connection: the server keeps the socket
// open but stops sending (the black-holed form where the OS never delivers
// FIN/RST). Once ReadKeepalive expires, receiveLoop must exit and trigger
// onDisconnect. Without a deadline, receiveLoop would block forever and the
// connection would go zombie (no events, no onDisconnect).
func TestReceiveLoop_ReadKeepaliveDeadLink(t *testing.T) {
	server, client := net.Pipe()

	disconnected := make(chan struct{})
	opts := DefaultInboundOptions
	opts.Password = "ClueCon"
	opts.ReadKeepalive = 200 * time.Millisecond
	opts.OnDisconnect = func() {
		close(disconnected)
	}

	connection := newConnection(client, false, opts.Options)
	go connection.disconnectLoop(opts.OnDisconnect)

	// Put the connection into its normal state
	_, err := server.Write([]byte("Content-Type: auth/request\r\n\r\n"))
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// The server goes silent: no close, no more bytes. After 200ms the read
	// keepalive expires, the read fails with a deadline exceeded, and the
	// connection is declared dead.

	select {
	case <-disconnected:
		// Expected: keepalive deadline detected the dead link, onDisconnect fired
	case <-time.After(3 * time.Second):
		t.Fatal("silently open connection was not declared dead within the read keepalive window (onDisconnect not triggered)")
	}

	connection.Close()
	server.Close()
	client.Close()
}

// TestReceiveLoop_ReadKeepaliveRefreshOnTraffic verifies that the deadline is
// refreshed after every received message: the server sends data every 100ms
// (each frame well within the 300ms keepalive window), so the connection must
// stay alive for the full second and never be killed spuriously.
func TestReceiveLoop_ReadKeepaliveRefreshOnTraffic(t *testing.T) {
	server, client := net.Pipe()

	opts := DefaultInboundOptions
	opts.Password = "ClueCon"
	opts.ReadKeepalive = 300 * time.Millisecond
	fired := make(chan struct{})
	opts.OnDisconnect = func() {
		select {
		case <-fired:
		default:
			close(fired)
		}
	}

	connection := newConnection(client, false, opts.Options)
	go connection.disconnectLoop(opts.OnDisconnect)

	_, err := server.Write([]byte("Content-Type: auth/request\r\n\r\n"))
	assert.NoError(t, err)

	// Sustained light traffic for 1s (one text/event-plain frame every 100ms,
	// each within the 300ms keepalive window). Event frames are used instead
	// of auth/request on purpose: eventLoop keeps consuming TypeEventPlain
	// (mirroring a production event stream), while auth/request is consumed
	// by authLoop, which is not started here -- once its channel buffer fills,
	// doMessage would block and pollute the timing.
	eventFrame := []byte("Content-Type: text/event-plain\r\nContent-Length: 25\r\n\r\nEvent-Name: HEARTBEAT\r\n\r\n")
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		_, err := server.Write(eventFrame)
		assert.NoError(t, err)
	}

	select {
	case <-fired:
		t.Fatal("read keepalive killed a connection with live traffic (deadline should have been refreshed)")
	default:
		// Expected: connection still alive
	}

	connection.Close()
	server.Close()
	client.Close()
}

// TestReadKeepaliveDisabledByDefault guards the opt-in contract: both default
// option sets leave ReadKeepalive at zero (disabled), so this feature never
// changes behavior for existing callers. Enabling the read deadline is an
// explicit, per-connection decision.
func TestReadKeepaliveDisabledByDefault(t *testing.T) {
	assert.Equal(t, time.Duration(0), DefaultOptions.ReadKeepalive,
		"generic defaults must keep the read keepalive disabled (opt-in)")
	assert.Equal(t, time.Duration(0), DefaultInboundOptions.ReadKeepalive,
		"inbound defaults must keep the read keepalive disabled (opt-in); enable it explicitly per connection")
}
