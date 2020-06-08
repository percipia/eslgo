package freeswitchesl

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

// Helper to check response status, only used for Auth checking at the moment
func (r RawResponse) IsOk() bool {
	return strings.HasPrefix(r.Headers.Get("Reply-Text"), "+OK")
}

// Helper to get the channel UUID
func (r RawResponse) ChannelUUID() string {
	return r.Headers.Get("Unique-ID")
}

// Helper function to get "Variable_" headers
func (r RawResponse) GetVariable(variable string) string {
	return r.GetHeader(fmt.Sprintf("Variable_%s", variable))
}

// Helper function that calls r.Header.Get
func (r RawResponse) GetHeader(header string) string {
	value, _ := url.PathUnescape(r.Headers.Get(header))
	return value
}

// Implement the Stringer interface for pretty printing
func (r RawResponse) String() string {
	var builder strings.Builder
	for key, values := range r.Headers {
		builder.WriteString(fmt.Sprintf("%s: %#v\n", key, values))
	}
	builder.Write(r.Body)
	return builder.String()
}

// Implement the GoStringer interface for pretty printing (%#v)
func (r RawResponse) GoString() string {
	return r.String()
}
