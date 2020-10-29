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
	"fmt"
	"io"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

const (
	TypeEventPlain  = `text/event-plain`
	TypeEventJSON   = `text/event-json`
	TypeEventXML    = `text/event-xml`
	TypeReply       = `command/reply`
	TypeAPIResponse = `api/response`
	TypeAuthRequest = `auth/request`
	TypeDisconnect  = `text/disconnect-notice`
)

// RawResponse This struct contains all response data from FreeSWITCH
type RawResponse struct {
	Headers textproto.MIMEHeader
	Body    []byte
}

func (c *Conn) readResponse() (*RawResponse, error) {
	header, err := c.header.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	response := &RawResponse{
		Headers: header,
	}

	if contentLength := header.Get("Content-Length"); len(contentLength) > 0 {
		length, err := strconv.Atoi(contentLength)
		if err != nil {
			return response, err
		}
		response.Body = make([]byte, length)
		_, err = io.ReadFull(c.reader, response.Body)
		if err != nil {
			return response, err
		}
	}

	return response, nil
}

// IsOk Helper to check response status, uses the Reply-Text header primarily. Calls GetReply internally
func (r RawResponse) IsOk() bool {
	return strings.HasPrefix(r.GetReply(), "+OK")
}

// GetReply Helper to get the Reply text from FreeSWITCH, uses the Reply-Text header primarily.
// Also will use the body if the Reply-Text header does not exist, this can be the case for TypeAPIResponse
func (r RawResponse) GetReply() string {
	if r.HasHeader("Reply-Text") {
		return r.GetHeader("Reply-Text")
	}
	return string(r.Body)
}

// ChannelUUID Helper to get the channel UUID. Calls GetHeader internally
func (r RawResponse) ChannelUUID() string {
	return r.GetHeader("Unique-ID")
}

// HasHeader Helper to check if the RawResponse has a header
func (r RawResponse) HasHeader(header string) bool {
	_, ok := r.Headers[textproto.CanonicalMIMEHeaderKey(header)]
	return ok
}

// GetVariable Helper function to get "Variable_" headers. Calls GetHeader internally
func (r RawResponse) GetVariable(variable string) string {
	return r.GetHeader(fmt.Sprintf("Variable_%s", variable))
}

// GetHeader Helper function that calls RawResponse.Headers.Get. Result gets passed through url.PathUnescape
func (r RawResponse) GetHeader(header string) string {
	value, _ := url.PathUnescape(r.Headers.Get(header))
	return value
}

// String Implement the Stringer interface for pretty printing
func (r RawResponse) String() string {
	var builder strings.Builder
	for key, values := range r.Headers {
		builder.WriteString(fmt.Sprintf("%s: %#v\n", key, values))
	}
	builder.Write(r.Body)
	return builder.String()
}

// GoString Implement the GoStringer interface for pretty printing (%#v)
func (r RawResponse) GoString() string {
	return r.String()
}
