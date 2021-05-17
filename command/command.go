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
	"net/textproto"
	"sort"
	"strings"
)

// Command - A basic interface for FreeSWITCH ESL commands. Implement this if you want to send your own raw data to FreeSIWTCH over the ESL connection. Do not add the eslgo.EndOfMessage(\r\n\r\n) marker, eslgo does that for you.
type Command interface {
	BuildMessage() string
}

var crlfToLF = strings.NewReplacer("\r\n", "\n")

// FormatHeaderString - Writes headers in a FreeSWITCH ESL friendly format. Converts headers containing \r\n to \n
func FormatHeaderString(headers textproto.MIMEHeader) string {
	var ws strings.Builder

	keys := make([]string, len(headers))
	i := 0
	for key := range headers {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	for _, key := range keys {
		for _, value := range headers[key] {
			value = crlfToLF.Replace(value)
			value = textproto.TrimString(value)
			ws.WriteString(key)
			ws.WriteString(": ")
			ws.WriteString(value)
			ws.WriteString("\r\n")
		}
	}
	// Remove the extra \r\n
	return ws.String()[:ws.Len()-2]
}
