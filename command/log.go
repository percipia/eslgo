package command

import "fmt"

type Log struct {
	Enabled bool
	Level   int
}

func (l Log) BuildMessage() string {
	if l.Enabled {
		return fmt.Sprintf("log %d", l.Level)
	}
	return "nolog"
}
