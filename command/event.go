package command

import (
	"fmt"
	"strings"
)

type Event struct {
	Format string
	Listen []string
}

func (e Event) String() string {
	return fmt.Sprintf("event %s, %s", e.Format, strings.Join(e.Listen, " "))
}
