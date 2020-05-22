package freeswitchesl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/textproto"
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

func readXMLEvent(body []byte) (*Event, error) {
	return &Event{}, nil
}

func readJSONEvent(body []byte) (*Event, error) {
	return &Event{}, nil
}

// Implement the Stringer interface for pretty printing (%v)
func (e Event) String() string {
	var builder strings.Builder
	for key, values := range e.Headers {
		builder.WriteString(fmt.Sprintf("%s: %#v\n", key, values))
	}
	builder.Write(e.Body)
	return builder.String()
}

// Implement the GoStringer interface for pretty printing (%#v)
func (e Event) GoString() string {
	return e.String()
}
