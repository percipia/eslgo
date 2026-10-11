/*
 * Copyright (c) 2020 Percipia
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Contributor(s):
 * Glenn O. Larsen <glenn.larsen@gmail.com>
 */
package eslgo

import (
	"bufio"
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInboundDial_ReturnsErrorWhenConnectionClosesBeforeAuthRequest(t *testing.T) {
	listener := listenTCP(t)
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		conn.Close()
	}()

	opts := DefaultInboundOptions
	opts.AuthTimeout = 100 * time.Millisecond
	opts.Logger = NilLogger{}
	disconnected := make(chan struct{})
	opts.OnDisconnect = func() {
		close(disconnected)
	}

	started := time.Now()
	conn, err := opts.Dial(listener.Addr().String())

	assert.Nil(t, conn)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "before auth request")
	assert.Less(t, time.Since(started), time.Second)
	assertOnDisconnect(t, disconnected)
}

func TestInboundDial_ReturnsErrorWhenAuthRequestTimesOut(t *testing.T) {
	listener := listenTCP(t)
	defer listener.Close()

	done := make(chan struct{})
	defer close(done)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		<-done
	}()

	opts := DefaultInboundOptions
	opts.AuthTimeout = 100 * time.Millisecond
	opts.Logger = NilLogger{}
	disconnected := make(chan struct{})
	opts.OnDisconnect = func() {
		close(disconnected)
	}

	started := time.Now()
	conn, err := opts.Dial(listener.Addr().String())

	assert.Nil(t, conn)
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
	assert.Less(t, time.Since(started), time.Second)
	assertOnDisconnect(t, disconnected)
}

func assertOnDisconnect(t *testing.T, disconnected <-chan struct{}) {
	t.Helper()
	select {
	case <-disconnected:
	case <-time.After(time.Second):
		t.Fatal("expected OnDisconnect to be called")
	}
}

func listenTCP(t *testing.T) net.Listener {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	return listener
}

// TestInboundDial_CloseWhileWaitingForAuthRequest exercises the first two
// guarded responseChannels read sites: the map lookups in Dial's pre-auth
// select. When the server hangs up before sending the auth request, the
// disconnect message must wake that select, Dial must fail with the
// "before auth request" error, and OnDisconnect must fire. This guards the
// behavior of those read sites under the responseChanMutex fix; the
// concurrent DATA RACE reproduction lives in
// TestInboundDial_CloseRightAfterSuccessfulAuth and
// TestOutbound_CloseWhileConnectionIsServed, where the closing side runs on
// a different goroutine than the readers.
func TestInboundDial_CloseWhileWaitingForAuthRequest(t *testing.T) {
	for i := 0; i < 25; i++ {
		listener := listenTCP(t)

		disconnected := make(chan struct{})
		var once sync.Once
		opts := DefaultInboundOptions
		opts.AuthTimeout = 2 * time.Second
		opts.Logger = NilLogger{}
		opts.OnDisconnect = func() {
			once.Do(func() { close(disconnected) })
		}

		go func() {
			serverConn, err := listener.Accept()
			if err != nil {
				return
			}
			// Give Dial time to block in its pre-auth select before hanging
			// up, so the disconnect message wakes that select instead of the
			// auth timeout.
			time.Sleep(10 * time.Millisecond)
			serverConn.Close()
		}()

		conn, err := opts.Dial(listener.Addr().String())

		assert.Nil(t, conn)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "before auth request")
		assertOnDisconnect(t, disconnected)
		listener.Close()
	}
}

// TestInboundDial_CloseRightAfterSuccessfulAuth hammers the window where the
// loops spawned by Dial (authLoop, disconnectLoop) evaluate their
// responseChannels map lookups while (*Conn).close deletes those entries.
// Run with -race: without the read-side RLocks this reliably reports a DATA
// RACE between the selects in the loops and the delete in close. Each
// iteration also proves the happy path still works: the connection
// authenticates against a loopback fake and is then torn down cleanly.
func TestInboundDial_CloseRightAfterSuccessfulAuth(t *testing.T) {
	for i := 0; i < 25; i++ {
		listener := listenTCP(t)

		var serverDone sync.WaitGroup
		serverDone.Add(1)
		go func() {
			defer serverDone.Done()
			serverConn, err := listener.Accept()
			if err != nil {
				return
			}
			defer serverConn.Close()
			if _, err := serverConn.Write([]byte("Content-Type: auth/request\r\n\r\n")); err != nil {
				return
			}
			reader := bufio.NewReader(serverConn)
			if _, err := reader.ReadString('\n'); err != nil {
				// Client gave up before answering the auth request.
				return
			}
			if _, err := serverConn.Write([]byte("Content-Type: command/reply\r\nReply-Text: +OK accepted\r\n\r\n")); err != nil {
				return
			}
			// Block until the client hangs up so every iteration starts from
			// a fully torn down connection.
			_, _ = reader.ReadString('\n')
		}()

		opts := DefaultInboundOptions
		opts.AuthTimeout = 2 * time.Second
		opts.Logger = NilLogger{}
		conn, err := opts.Dial(listener.Addr().String())
		require.NoError(t, err)
		require.NotNil(t, conn)

		// Park the loop goroutines in their select statements, then close the
		// connection: the map reads inside those selects must not race with
		// the map deletes performed by (*Conn).close.
		time.Sleep(10 * time.Millisecond)
		conn.Close()

		serverDone.Wait()
		listener.Close()
		time.Sleep(5 * time.Millisecond)
	}
}
