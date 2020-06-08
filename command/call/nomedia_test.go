package call

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var TestNoMediaMessage = strings.ReplaceAll(`sendmsg none
Call-Command: nomedia
Nomedia-Uuid: test`, "\n", "\r\n")

func TestNoMedia_BuildMessage(t *testing.T) {
	nomedia := NoMedia{
		UUID:        "none",
		NoMediaUUID: "test",
	}
	assert.Equal(t, TestNoMediaMessage, nomedia.BuildMessage())
}
