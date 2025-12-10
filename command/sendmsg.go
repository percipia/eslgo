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
	"fmt"
	"net/textproto"
	"strconv"
	"strings"
)

type SendMessage struct {
	UUID    string
	Headers textproto.MIMEHeader
	Body    string
	Sync    bool
	SyncPri bool
}

func (s *SendMessage) BuildMessage() string {
	var headers []string

	if s.Headers == nil {
		s.Headers = make(textproto.MIMEHeader)
	}

	// Waits for this event to finish before continuing even in async mode
	if s.Sync {
		headers = append(headers, "event-lock: true")
	}

	// No documentation on this flag, I assume it takes priority over the other flag?
	if s.SyncPri {
		headers = append(headers, "event-lock-pri: true")
	}

	// Ensure the correct content length is set in the header
	if len(s.Body) > 0 {
		headers = append(headers, "Content-Length: "+strconv.Itoa(len(s.Body)))
	} else {
		delete(s.Headers, "Content-Length")
	}

	// Format the headers
	for key, values := range s.Headers {
		for _, value := range values {
			headers = append(headers, key+": "+value)
		}
	}

	headerString := strings.Join(headers, "\r\n")

	if len(s.Body) > 0 {
		return fmt.Sprintf("sendmsg %s\r\n%s\r\n\r\n%s", s.UUID, headerString, s.Body)
	}
	return fmt.Sprintf("sendmsg %s\r\n%s", s.UUID, headerString)
}
