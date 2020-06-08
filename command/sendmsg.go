package command

import (
	"fmt"
	"net/http"
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
	err := http.Header(s.Headers).Write(&headers)
	if err != nil || headers.Len() < 3 {
		return ""
	}
	// -2 to remove the trailing \r\n added by http.Header.Write
	headerString := headers.String()[:headers.Len()-2]
	if _, ok := s.Headers["Content-Length"]; ok {
		return fmt.Sprintf("sendmsg %s\r\n%s\r\n\r\n%s", s.UUID, headerString, s.Body)
	}
	return fmt.Sprintf("sendmsg %s\r\n%s", s.UUID, headerString)
}
