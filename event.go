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
package eslgo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

type EventListener func(event *Event)

type Event struct {
	Headers textproto.MIMEHeader
	Body    []byte
}

const (
	EventListenAll = "ALL"
)

func readPlainEvent(body []byte) (*Event, error) {
	reader := bufio.NewReader(bytes.NewBuffer(body))
	header := textproto.NewReader(reader)

	headers, err := header.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	event := &Event{
		Headers: headers,
	}

	if contentLength := headers.Get("Content-Length"); len(contentLength) > 0 {
		length, err := strconv.Atoi(contentLength)
		if err != nil {
			return event, err
		}
		event.Body = make([]byte, length)
		_, err = io.ReadFull(reader, event.Body)
		if err != nil {
			return event, err
		}
	}

	return event, nil
}

// TODO: Needs processing
func readXMLEvent(body []byte) (*Event, error) {
	return &Event{
		Headers: make(textproto.MIMEHeader),
	}, nil
}

// TODO: Needs processing
func readJSONEvent(body []byte) (*Event, error) {
	return &Event{
		Headers: make(textproto.MIMEHeader),
	}, nil
}

// GetName Helper function that returns the event name header
func (e Event) GetName() string {
	return e.GetHeader("Event-Name")
}

// HasHeader Helper to check if the Event has a header
func (e Event) HasHeader(header string) bool {
	_, ok := e.Headers[textproto.CanonicalMIMEHeaderKey(header)]
	return ok
}

// GetHeader Helper function that calls e.Header.Get
func (e Event) GetHeader(header string) string {
	value, _ := url.PathUnescape(e.Headers.Get(header))
	return value
}

// String Implement the Stringer interface for pretty printing (%v)
func (e Event) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s\n", e.GetName()))
	for key, values := range e.Headers {
		builder.WriteString(fmt.Sprintf("%s: %#v\n", key, values))
	}
	builder.Write(e.Body)
	return builder.String()
}

// GoString Implement the GoStringer interface for pretty printing (%#v)
func (e Event) GoString() string {
	return e.String()
}
