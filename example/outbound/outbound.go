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
	"context"
	"fmt"
	"github.com/percipia/eslgo"
	"log"
)

func main() {
	// Start listening, this is a blocking function
	log.Fatalln(eslgo.ListenAndServe(":8084", handleConnection))
}

func handleConnection(ctx context.Context, conn *eslgo.Conn, response *eslgo.RawResponse) {
	fmt.Printf("Got connection! %#v\n", response)

	// Place the call in the foreground(api) to user 100 and playback an audio file as the bLeg and no exported variables
	response, err := conn.OriginateCall(ctx, false, eslgo.Leg{CallURL: "user/100"}, eslgo.Leg{CallURL: "&playback(misc/ivr-to_hear_screaming_monkeys.wav)"}, map[string]string{})
	fmt.Println("Call Originated: ", response, err)
}
