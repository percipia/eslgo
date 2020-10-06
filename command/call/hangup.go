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
	"github.com/percipia/eslgo/command"
	"net/textproto"
)

type Hangup struct {
	UUID    string
	Cause   string
	Sync    bool
	SyncPri bool
}

func (h Hangup) BuildMessage() string {
	sendMsg := command.SendMessage{
		UUID:    h.UUID,
		Headers: make(textproto.MIMEHeader),
		Sync:    h.Sync,
		SyncPri: h.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "hangup")
	sendMsg.Headers.Set("hangup-cause", h.Cause)

	return sendMsg.BuildMessage()
}
