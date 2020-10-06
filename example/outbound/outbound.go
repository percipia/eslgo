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
	log.Fatalln(eslgo.ListenAndServe(":8084", handleConnection))
}

func handleConnection(ctx context.Context, conn *eslgo.Conn, response *eslgo.RawResponse) {
	fmt.Printf("Got connection! %#v\n", response)
	_ = conn.EnableEvents(ctx)
	originationUUID, response, err := conn.OriginateCall(ctx, "user/100", "&playback(misc/ivr-to_hear_screaming_monkeys.wav)", map[string]string{})
	fmt.Println("Call Originated: ", originationUUID, response, err)
}
