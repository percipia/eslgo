package call

import (
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"net"
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
		Sync:    u.Sync,
		SyncPri: u.SyncPri,
	}
	localHost, localPort, _ := net.SplitHostPort(u.Local.String())
	remoteHost, remotePort, _ := net.SplitHostPort(u.Local.String())
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
