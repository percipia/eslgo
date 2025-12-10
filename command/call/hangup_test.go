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

var TestHangupMessage = strings.ReplaceAll(`sendmsg none
Call-Command: hangup
Hangup-Cause: NORMAL_CLEARING`, "\n", "\r\n")

func TestHangup_BuildMessage(t *testing.T) {
	hangup := Hangup{
		UUID:  "none",
		Cause: "NORMAL_CLEARING",
	}
	assert.Equal(t, normalizeMessage(TestHangupMessage), normalizeMessage(hangup.BuildMessage()))
}
