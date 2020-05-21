package command

import "fmt"

type API struct {
	Command   string
	Arguments string
}

func (api API) String() string {
	return fmt.Sprintf("api %s %s", api.Command, api.Arguments)
}