package call

import (
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
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
