package call

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Normalizes the headers in a message to ensure they are always in the same order before comparison
func normalizeMessage(message string) string {
	parts := strings.Split(message, "\r\n\r\n")
	headers := strings.Split(parts[0], "\r\n")
	sort.Strings(headers)
	normalizedHeaders := strings.Join(headers, "\r\n")

	if len(parts) > 1 {
		return normalizedHeaders + "\r\n\r\n" + parts[1]
	}
	return normalizedHeaders
}

var (
	TestExecMessage = strings.ReplaceAll(`sendmsg none
Call-Command: execute
Execute-App-Arg: /tmp/test.wav
Execute-App-Name: playback
Loops: 1`, "\n", "\r\n")
	TestSetMessage = strings.ReplaceAll(`sendmsg none
Call-Command: execute
Content-Length: 11
Content-Type: text/plain
Execute-App-Name: set
Loops: 1

hello=world`, "\n", "\r\n")
	TestExportMessage = strings.ReplaceAll(`sendmsg none
Call-Command: execute
Content-Length: 11
Content-Type: text/plain
Execute-App-Name: export
Loops: 1

hello=world`, "\n", "\r\n")
	TestPushMessage = strings.ReplaceAll(`sendmsg none
Call-Command: execute
Content-Length: 11
Content-Type: text/plain
Execute-App-Name: push
Loops: 1

hello=world`, "\n", "\r\n")
)

func TestExecute_BuildMessage(t *testing.T) {
	exec := Execute{
		UUID:    "none",
		AppName: "playback",
		AppArgs: "/tmp/test.wav",
	}
	assert.Equal(t, normalizeMessage(TestExecMessage), normalizeMessage(exec.BuildMessage()))
}

func TestSet_BuildMessage(t *testing.T) {
	set := Set{
		UUID:  "none",
		Key:   "hello",
		Value: "world",
	}
	assert.Equal(t, normalizeMessage(TestSetMessage), normalizeMessage(set.BuildMessage()))
}

func TestExport_BuildMessage(t *testing.T) {
	export := Export{
		UUID:  "none",
		Key:   "hello",
		Value: "world",
	}
	assert.Equal(t, normalizeMessage(TestExportMessage), normalizeMessage(export.BuildMessage()))
}

func TestPush_BuildMessage(t *testing.T) {
	push := Push{
		UUID:  "none",
		Key:   "hello",
		Value: "world",
	}
	assert.Equal(t, normalizeMessage(TestPushMessage), normalizeMessage(push.BuildMessage()))
}
