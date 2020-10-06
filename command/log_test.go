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

func TestLog_BuildMessage(t *testing.T) {
	assert.Equal(t, "log 9", Log{Enabled: true, Level: 9}.BuildMessage())
}

func TestNoLog_BuildMessage(t *testing.T) {
	assert.Equal(t, "nolog", Log{}.BuildMessage())
}
