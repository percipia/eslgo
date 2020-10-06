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
package call

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	TestExecMessage = strings.ReplaceAll(`sendmsg none
Call-Command: execute
Execute-App-Arg: /tmp/test.wav
Execute-App-Name: playback
Loops: 1`, "\n", "\r\n")
	TestSetMessage = strings.ReplaceAll(`sendmsg none
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
	assert.Equal(t, TestExecMessage, exec.BuildMessage())
}

func TestSet_BuildMessage(t *testing.T) {
	set := Set{
		UUID:  "none",
		Key:   "hello",
		Value: "world",
	}
	assert.Equal(t, TestSetMessage, set.BuildMessage())
}
