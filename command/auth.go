package command

import "fmt"

type Auth struct {
	User     string
	Password string
}

func (auth Auth) String() string {
	if len(auth.User) > 0 {
		return fmt.Sprintf("userauth %s:%s", auth.User, auth.Password)
	}
	return fmt.Sprintf("auth %s", auth.Password)
}