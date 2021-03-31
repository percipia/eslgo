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
