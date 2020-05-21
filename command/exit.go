package command

type Exit struct {}

func (Exit) String() string {
	return "exit"
}