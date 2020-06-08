package call

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var TestHangupMessage = strings.ReplaceAll(`sendmsg none
Call-Command: hangup
Hangup-Cause: NORMAL_CLEARING`, "\n", "\r\n")

func TestHangup_BuildMessage(t *testing.T) {
	hangup := Hangup{
		UUID:  "none",
		Cause: "NORMAL_CLEARING",
	}
	assert.Equal(t, TestHangupMessage, hangup.BuildMessage())
}
