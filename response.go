package freeswitchesl

import (
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

const (
	TypeEventPlain  = `text/event-plain`
	TypeEventJSON   = `text/event-json`
	TypeEventXML    = `text/event-xml`
	TypeReply       = `command/reply`
	TypeAPIResponse = `api/response`
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

func (r RawResponse) IsOk() bool {
	return strings.HasPrefix(r.Headers.Get("Reply-Text"), "+OK")
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
