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

type Event struct {
	Ignore bool
	Format string
	Listen []string
}

type MyEvents struct {
	Format string
	UUID   string
}

type DisableEvents struct{}

// The divert_events command is available to allow events that an embedded script would expect to get in the inputcallback to be diverted to the event socket.
type DivertEvents struct {
	Enabled bool
}

type SendEvent struct {
	Name    string
	Headers textproto.MIMEHeader
	Body    string
}

func (e Event) BuildMessage() string {
	prefix := ""
	if e.Ignore {
		prefix = "nix"
	}
	return fmt.Sprintf("%sevent %s %s", prefix, e.Format, strings.Join(e.Listen, " "))
}

func (m MyEvents) BuildMessage() string {
	if len(m.UUID) > 0 {
		return fmt.Sprintf("myevents %s %s", m.Format, m.UUID)

	}
	return fmt.Sprintf("myevents %s", m.Format)
}

func (DisableEvents) BuildMessage() string {
	return "noevents"
}

func (d DivertEvents) BuildMessage() string {
	if d.Enabled {
		return "divert_events on"
	}
	return "divert_events off"
}

func (s *SendEvent) BuildMessage() string {
	// Ensure the correct content length is set in the header
	if len(s.Body) > 0 {
		s.Headers.Set("Content-Length", strconv.Itoa(len(s.Body)))
	} else {
		delete(s.Headers, "Content-Length")
	}

	// Format the headers
	headerString := FormatHeaderString(s.Headers)
	if _, ok := s.Headers["Content-Length"]; ok {
		return fmt.Sprintf("sendevent %s\r\n%s\r\n\r\n%s", s.Name, headerString, s.Body)
	}
	return fmt.Sprintf("sendevent %s\r\n%s", s.Name, headerString)
}
