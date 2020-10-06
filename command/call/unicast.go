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
	"net"
	"net/textproto"
)

/*
 * unicast is used to hook up mod_spandsp for faxing over a socket.
 * Note:
 * That is a nice way for a script or app that uses the socket interface to get at the media.
 * It's good because then spandsp isn't living inside of FreeSWITCH and it can run on a box sitting next to it. It scales better.
 */
type Unicast struct {
	UUID    string
	Local   net.Addr
	Remote  net.Addr
	Flags   string
	Sync    bool
	SyncPri bool
}

func (u Unicast) BuildMessage() string {
	sendMsg := command.SendMessage{
		UUID:    u.UUID,
		Headers: make(textproto.MIMEHeader),
		Sync:    u.Sync,
		SyncPri: u.SyncPri,
	}
	localHost, localPort, _ := net.SplitHostPort(u.Local.String())
	remoteHost, remotePort, _ := net.SplitHostPort(u.Remote.String())
	sendMsg.Headers.Set("call-command", "unicast")
	sendMsg.Headers.Set("local-ip", localHost)
	sendMsg.Headers.Set("local-port", localPort)
	sendMsg.Headers.Set("remote-ip", remoteHost)
	sendMsg.Headers.Set("remote-port", remotePort)
	sendMsg.Headers.Set("transport", u.Local.Network())
	if len(u.Flags) > 0 {
		sendMsg.Headers.Set("flags", u.Flags)
	}
	return sendMsg.BuildMessage()
}
