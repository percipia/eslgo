/*
 * Copyright (c) 2026 Percipia
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Contributor(s):
 * anrenlx2025
 */
package eslgo

import (
	"bufio"
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOutbound_CloseWhileConnectionIsServed(t *testing.T) {
	for i := 0; i < 25; i++ {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)

		served := make(chan *Conn, 1)
		handler := func(ctx context.Context, conn *Conn, connectResponse *RawResponse) {
			served <- conn
		}

		opts := DefaultOutboundOptions
		opts.Logger = NilLogger{}
		accepted := make(chan struct{})
		go func() {
			serverConn, err := listener.Accept()
			close(accepted)
			if err != nil {
				return
			}
			// No Close on serverConn here: the eslgo Conn owns the socket for
			// the rest of the iteration, exactly like ListenAndServe.
			conn := newConnection(serverConn, true, opts.Options)
			go conn.dummyLoop()
			go conn.outboundHandle(handler, opts.ConnectionDelay, opts.ConnectTimeout)
		}()

		fakeFS, err := net.Dial("tcp", listener.Addr().String())
		require.NoError(t, err)
		<-accepted

		// Complete the outbound handshake so the handler runs while
		// dummyLoop is parked in its select over responseChannels.
		reader := bufio.NewReader(fakeFS)
		_, err = reader.ReadString('\n')
		require.NoError(t, err)
		_, err = fakeFS.Write([]byte("Content-Type: command/reply\r\nReply-Text: +OK\r\n\r\n"))
		require.NoError(t, err)

		conn := <-served
		require.NotNil(t, conn)

		// Closing here makes (*Conn).close delete the responseChannels map
		// entries while dummyLoop reads them inside its select.
		time.Sleep(10 * time.Millisecond)
		conn.Close()

		_ = fakeFS.Close()
		listener.Close()
		time.Sleep(5 * time.Millisecond)
	}
}
