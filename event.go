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
	"encoding/json"
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

	// Use the lenient ESL header parser instead of ReadMIMEHeader. FreeSWITCH
	// permits characters in variable names (for example "variable_fish,(sticks)~_1")
	// that net/textproto rejects, which previously caused the entire event to be
	// discarded with a MIME parsing error.
	headers, err := readESLMIMEHeader(header)
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

/*
 * readJSONEvent parses a text/event-json payload into the same Event shape the
 * plaintext parser produces, so callers get Headers populated instead of having
 * to unmarshal the body themselves.
 *
 * FreeSWITCH emits the event as a flat JSON object of header name to value. The
 * special "_body" key carries the event body when one is present, matching the
 * Content-Length body of the plaintext format. Header names are canonicalized
 * with the same rules as the plaintext path, so Get/HasHeader behave identically
 * regardless of the negotiated event format.
 */
func readJSONEvent(body []byte) (*Event, error) {
	event := &Event{
		Headers: make(textproto.MIMEHeader),
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		// Keep the raw payload so nothing is silently lost if FreeSWITCH sends
		// something we cannot decode.
		event.Body = body
		return event, err
	}

	for key, raw := range decoded {
		value, ok := jsonHeaderValue(raw)
		if !ok {
			continue
		}
		if key == "_body" {
			event.Body = []byte(value)
			continue
		}
		canonical := canonicalESLHeaderKey(key)
		event.Headers[canonical] = append(event.Headers[canonical], value)
	}

	return event, nil
}

// jsonHeaderValue converts a decoded JSON value into its ESL header string form.
// Objects and arrays are skipped since FreeSWITCH event headers are always scalar.
func jsonHeaderValue(raw interface{}) (string, bool) {
	switch v := raw.(type) {
	case string:
		return v, true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(v), true
	case nil:
		return "", true
	default:
		return "", false
	}
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
