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

// Documentation is sparse on this, but it looks like it transfers a call to an application?
type Transfer struct {
	UUID        string
	Application string
	Sync        bool
	SyncPri     bool
}

func (t Transfer) BuildMessage() string {
	sendMsg := command.SendMessage{
		UUID:    t.UUID,
		Headers: make(textproto.MIMEHeader),
		Sync:    t.Sync,
		SyncPri: t.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "xferext")
	sendMsg.Headers.Set("application", t.Application)

	return sendMsg.BuildMessage()
}
