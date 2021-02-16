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
package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/percipia/eslgo"
	"os"
	"time"
)

func main() {
	// Connect to FreeSWITCH
	conn, err := eslgo.Dial("127.0.0.1:8021", "ClueCon", func() {
		fmt.Println("Inbound Connection Disconnected")
	})
	if err != nil {
		fmt.Println("Error connecting", err)
		return
	}

	// Register an event listener for all events
	listenerID := conn.RegisterEventListener(eslgo.EventListenAll, func(event *eslgo.Event) {
		fmt.Printf("%#v\n", event)
	})

	// Ensure all events are enabled
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = conn.EnableEvents(ctx)
	cancel()

	// Wait until enter is pressed to exit
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text != "" {
			break
		}
	}

	// Remove the listener and close the connection gracefully
	conn.RemoveEventListener(eslgo.EventListenAll, listenerID)
	conn.ExitAndClose()
}
