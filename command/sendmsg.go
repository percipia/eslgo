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
	if s.Headers == nil {
		s.Headers = make(textproto.MIMEHeader)
	}
	// Waits for this event to finish before continuing even in async mode
	if s.Sync {
		s.Headers.Set("event-lock", "true")
	}
	// No documentation on this flag, I assume it takes priority over the other flag?
	if s.SyncPri {
		s.Headers.Set("event-lock-pri", "true")
	}

	// Ensure the correct content length is set in the header
	if len(s.Body) > 0 {
		s.Headers.Set("Content-Length", strconv.Itoa(len(s.Body)))
	} else {
		delete(s.Headers, "Content-Length")
	}

	// Format the headers
	var headers strings.Builder
	for key, values := range s.Headers {
		for _, value := range values {
			headers.WriteString(key)
			headers.WriteString(": ")
			headers.WriteString(value)
		}
	}
	if _, ok := s.Headers["Content-Length"]; ok {
		return fmt.Sprintf("sendmsg %s\r\n%s\r\n\r\n%s", s.UUID, headers.String(), s.Body)
	}
	return fmt.Sprintf("sendmsg %s\r\n%s", s.UUID, headers.String())
}
