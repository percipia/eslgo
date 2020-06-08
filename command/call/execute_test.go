package call

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	ExecMessage = strings.ReplaceAll(`sendmsg none
Call-Command: execute
Execute-App-Arg: /tmp/test.wav
Execute-App-Name: playback
Loops: 1`, "\n", "\r\n")
	SetMessage = strings.ReplaceAll(`sendmsg none
Call-Command: execute
Execute-App-Arg: hello=world
Execute-App-Name: set
Loops: 1`, "\n", "\r\n")
)

func TestExecute_BuildMessage(t *testing.T) {
	exec := Execute{
		UUID:    "none",
		AppName: "playback",
		AppArgs: "/tmp/test.wav",
	}
	assert.Equal(t, ExecMessage, exec.BuildMessage())
}

func TestSet_BuildMessage(t *testing.T) {
	set := Set{
		UUID:  "none",
		Key:   "hello",
		Value: "world",
	}
	assert.Equal(t, SetMessage, set.BuildMessage())
}
