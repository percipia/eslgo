package call

import (
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
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
		Sync:    n.Sync,
		SyncPri: n.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "nomedia")
	sendMsg.Headers.Set("nomedia-uuid", n.NoMediaUUID)

	return sendMsg.BuildMessage()
}
