package command

import "fmt"

type API struct {
	Command    string
	Arguments  string
	Background bool
}

func (api API) String() string {
	if api.Background {
		return fmt.Sprintf("bgapi %s %s", api.Command, api.Arguments)
	}
	return fmt.Sprintf("api %s %s", api.Command, api.Arguments)
}
