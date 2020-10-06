/*
 * Copyright (c) 2020 Percipia
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Contributor(s):
 * Andrew Querol <aquerol@percipia.com>
 */
package command

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var TestSendEventMessage = strings.ReplaceAll(`sendevent MESSAGE_WAITING
MWI-Message-Account: 7100@192.168.1.1
MWI-Messages-Waiting: yes
MWI-Voice-Message: 5/5 (1/1)`, "\n", "\r\n")

func TestDisableEvents_BuildMessage(t *testing.T) {
	assert.Equal(t, "noevents", DisableEvents{}.BuildMessage())
}

func TestDivertEvents_BuildMessage(t *testing.T) {
	assert.Equal(t, "divert_events on", DivertEvents{true}.BuildMessage())
	assert.Equal(t, "divert_events off", DivertEvents{false}.BuildMessage())
}

func TestEvent_BuildMessage(t *testing.T) {
	assert.Equal(t, "event plain MESSAGE_QUERY", Event{
		Format: "plain",
		Listen: []string{"MESSAGE_QUERY"},
	}.BuildMessage())
	assert.Equal(t, "nixevent plain MESSAGE_QUERY", Event{
		Ignore: true,
		Format: "plain",
		Listen: []string{"MESSAGE_QUERY"},
	}.BuildMessage())
}

func TestMyEvents_BuildMessage(t *testing.T) {
	assert.Equal(t, "myevents plain none", MyEvents{Format: "plain", UUID: "none"}.BuildMessage())
}

func TestSendEvent_BuildMessage(t *testing.T) {
	sendEvent := SendEvent{
		Name: "MESSAGE_WAITING",
		Headers: map[string][]string{
			"MWI-Messages-Waiting": {"yes"},
			"MWI-Message-Account":  {"7100@192.168.1.1"},
			"MWI-Voice-Message":    {"5/5 (1/1)"},
		},
	}
	assert.Equal(t, TestSendEventMessage, sendEvent.BuildMessage())
}
