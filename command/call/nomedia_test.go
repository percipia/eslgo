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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestNoMediaMessage = strings.ReplaceAll(`sendmsg none
Call-Command: nomedia
Nomedia-Uuid: test`, "\n", "\r\n")

func TestNoMedia_BuildMessage(t *testing.T) {
	nomedia := NoMedia{
		UUID:        "none",
		NoMediaUUID: "test",
	}
	assert.Equal(t, normalizeMessage(TestNoMediaMessage), normalizeMessage(nomedia.BuildMessage()))
}
