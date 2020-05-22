package call

import (
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
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
		Sync:    h.Sync,
		SyncPri: h.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "hangup")
	sendMsg.Headers.Set("hangup-cause", h.Cause)

	return sendMsg.BuildMessage()
}
