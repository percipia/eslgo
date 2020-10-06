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
	"gitlab.percipia.com/libs/go/freeswitchesl"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"log"
	"time"
)

func main() {
	log.Fatalln(freeswitchesl.ListenAndServe(":8084", handleConnection))
}

func handleConnection(ctx context.Context, conn *freeswitchesl.Conn, response *freeswitchesl.RawResponse) {
	fmt.Printf("Got connection! %#v\n", response)
	conn.SendCommand(ctx, command.Event{
		Format: "plain",
		Listen: []string{"ALL"},
	})
	conn.SendCommand(ctx, command.API{
		Command:   "originate",
		Arguments: "user/100 &park()",
	})
	time.Sleep(60 * time.Second)
}
