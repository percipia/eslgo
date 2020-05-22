package call

import (
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
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
		Sync:    t.Sync,
		SyncPri: t.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "xferext")
	sendMsg.Headers.Set("application", t.Application)

	return sendMsg.BuildMessage()
}
