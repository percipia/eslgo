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
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestUnicastMessage = strings.ReplaceAll(`sendmsg none
Call-Command: unicast
Flags: native
Local-Ip: 192.168.1.100
Local-Port: 8025
Remote-Ip: 192.168.1.101
Remote-Port: 8026
Transport: tcp`, "\n", "\r\n")

func TestUnicast_BuildMessage(t *testing.T) {
	testLocal, _ := net.ResolveTCPAddr("tcp", "192.168.1.100:8025")
	testRemote, _ := net.ResolveTCPAddr("tcp", "192.168.1.101:8026")
	unicast := Unicast{
		UUID:   "none",
		Local:  testLocal,
		Remote: testRemote,
		Flags:  "native",
	}
	assert.Equal(t, normalizeMessage(TestUnicastMessage), normalizeMessage(unicast.BuildMessage()))
}
