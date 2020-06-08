package command

type Exit struct{}

func (Exit) BuildMessage() string {
	return "exit"
}
