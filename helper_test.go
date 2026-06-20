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
	"net"
	"net/textproto"
	"sync"
	"testing"
)

// TestWaitForDTMF_NoPanicOnLateEvent reproduces the "send on closed channel"
// panic from WaitForDTMF: a DTMF listener goroutine is dispatched by
// callEventListener and may still be running after WaitForDTMF has returned
// (on ctx cancellation) and closed its done channel.
func TestWaitForDTMF_NoPanicOnLateEvent(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	conn := newConnection(client, false, DefaultOptions)
	defer conn.Close()

	const uuid = "test-uuid"

	dtmfEvent := func() *Event {
		headers := make(textproto.MIMEHeader)
		headers.Set("Event-Name", "DTMF")
		headers.Set("Unique-Id", uuid)
		headers.Set("DTMF-Digit", "1")
		return &Event{Headers: headers}
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Fire DTMF events concurrently so a listener goroutine is
			// in-flight when WaitForDTMF returns and closes its channel.
			for j := 0; j < 5; j++ {
				conn.callEventListener(dtmfEvent())
			}
		}()
		cancel()
		_, _ = conn.WaitForDTMF(ctx, uuid)
	}
	wg.Wait()
}
