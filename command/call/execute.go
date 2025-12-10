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
package call

import (
	"fmt"
	"net/textproto"
	"strconv"

	"github.com/percipia/eslgo/command"
)

type Execute struct {
	UUID      string
	AppName   string
	AppArgs   string
	AppUUID   string
	Loops     int
	Sync      bool
	SyncPri   bool
	ForceBody bool
}

// Helper to call Execute with Set since it is commonly used
type Set struct {
	UUID    string
	Key     string
	Value   string
	Sync    bool
	SyncPri bool
}

// Helper to call Execute with Export since it is commonly used
type Export Set

// Helper to call Execute with Push since it is commonly used
type Push Set

func (s Set) buildMessage(app string) string {
	e := Execute{
		UUID:      s.UUID,
		AppName:   app,
		AppArgs:   fmt.Sprintf("%s=%s", s.Key, s.Value),
		Sync:      s.Sync,
		SyncPri:   s.SyncPri,
		ForceBody: true,
	}
	return e.BuildMessage()
}

func (s Set) BuildMessage() string {
	return s.buildMessage("set")
}

func (e Export) BuildMessage() string {
	return Set(e).buildMessage("export")
}

func (p Push) BuildMessage() string {
	return Set(p).buildMessage("push")
}

func (e *Execute) BuildMessage() string {
	if e.Loops == 0 {
		e.Loops = 1
	}
	sendMsg := command.SendMessage{
		UUID:    e.UUID,
		Headers: make(textproto.MIMEHeader, 4), // preallocating for the 4+ headers that are always set to reduce amount of dynamic allocations
		Sync:    e.Sync,
		SyncPri: e.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "execute")
	sendMsg.Headers.Set("execute-app-name", e.AppName)
	sendMsg.Headers.Set("loops", strconv.Itoa(e.Loops))
	// This allows us to track when application execution completes via the Application-UUID header in events.
	if e.AppUUID != "" {
		sendMsg.Headers.Set("Event-UUID", e.AppUUID)
	}

	// According to documentation that is the max header length
	if len(e.AppArgs) > 2048 || e.ForceBody {
		sendMsg.Headers.Set("Content-Type", "text/plain")
		sendMsg.Body = e.AppArgs
	} else {
		sendMsg.Headers.Set("execute-app-arg", e.AppArgs)
	}

	return sendMsg.BuildMessage()
}
