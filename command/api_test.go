package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestAPIMessage   = `api originate user/100 &park()`
	TestBGAPIMessage = `bgapi originate user/100 &park()`
)

func TestAPI_BuildMessage(t *testing.T) {
	api := API{
		Command:   "originate",
		Arguments: "user/100 &park()",
	}
	assert.Equal(t, TestAPIMessage, api.BuildMessage())
}

func TestAPI_BuildMessage_BG(t *testing.T) {
	api := API{
		Command:    "originate",
		Arguments:  "user/100 &park()",
		Background: true,
	}
	assert.Equal(t, TestBGAPIMessage, api.BuildMessage())
}
