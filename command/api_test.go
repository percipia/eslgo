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
