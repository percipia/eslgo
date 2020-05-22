package call

import (
	"fmt"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"net/textproto"
	"strconv"
)

type Execute struct {
	UUID    string
	AppName string
	AppArgs string
	Loops   int
	Sync    bool
	SyncPri bool
}

// Helper to call Execute with Set or Export since it is commonly used
type Set struct {
	Export  bool
	UUID    string
	Key     string
	Value   string
	Sync    bool
	SyncPri bool
}

func (s Set) BuildMessage() string {
	e := Execute{
		UUID:    s.UUID,
		AppName: "set",
		AppArgs: fmt.Sprintf("%s=%s", s.Key, s.Value),
		Sync:    s.Sync,
		SyncPri: s.SyncPri,
	}
	if s.Export {
		e.AppName = "export"
	}
	return e.BuildMessage()
}

func (e *Execute) BuildMessage() string {
	if e.Loops == 0 {
		e.Loops = 1
	}
	sendMsg := command.SendMessage{
		UUID:    e.UUID,
		Headers: make(textproto.MIMEHeader),
		Sync:    e.Sync,
		SyncPri: e.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "execute")
	sendMsg.Headers.Set("execute-app-name", e.AppName)
	sendMsg.Headers.Set("loops", strconv.Itoa(e.Loops))

	// According to documentation that is the max header length
	if len(e.AppArgs) > 2048 {
		sendMsg.Headers.Set("content-type", "text/plain")
		sendMsg.Headers.Set("content-length", strconv.Itoa(len(e.AppArgs)))
		sendMsg.Body = e.AppArgs
	} else {
		sendMsg.Headers.Set("execute-app-arg", e.AppArgs)
	}

	return sendMsg.BuildMessage()
}
