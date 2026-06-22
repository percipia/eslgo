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
	"context"
	"errors"
	"net"
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
