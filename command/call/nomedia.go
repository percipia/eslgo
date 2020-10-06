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

type NoMedia struct {
	UUID        string
	NoMediaUUID string
	Sync        bool
	SyncPri     bool
}

func (n NoMedia) BuildMessage() string {
	sendMsg := command.SendMessage{
		UUID:    n.UUID,
		Headers: make(textproto.MIMEHeader),
		Sync:    n.Sync,
		SyncPri: n.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "nomedia")
	sendMsg.Headers.Set("nomedia-uuid", n.NoMediaUUID)

	return sendMsg.BuildMessage()
}
