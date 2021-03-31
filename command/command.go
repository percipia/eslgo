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

// Command - A basic interface for FreeSWITCH ESL commands. Implement this if you want to send your own raw data to FreeSIWTCH over the ESL connection. Do not add the eslgo.EndOfMessage(\r\n\r\n) marker, eslgo does that for you.
type Command interface {
	BuildMessage() string
}
