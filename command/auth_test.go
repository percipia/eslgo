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
	TestAuthMessage     = `auth testing123`
	TestUserAuthMessage = `userauth testuser:testing123`
)

func TestAuth_BuildMessage(t *testing.T) {
	auth := Auth{
		Password: "testing123",
	}
	assert.Equal(t, TestAuthMessage, auth.BuildMessage())
}

func TestAuth_BuildMessage_User(t *testing.T) {
	auth := Auth{
		User:     "testuser",
		Password: "testing123",
	}
	assert.Equal(t, TestUserAuthMessage, auth.BuildMessage())
}
