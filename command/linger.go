package command

type Linger struct {
	Enabled bool
}

func (l Linger) BuildMessage() string {
	if l.Enabled {
		return "linger"
	}
	return "nolinger"
}
