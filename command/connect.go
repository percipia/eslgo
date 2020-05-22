package command

type Connect struct{}

func (Connect) BuildMessage() string {
	return "connect"
}
