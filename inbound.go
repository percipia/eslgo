package freeswitchesl

import (
	"context"
	"errors"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"net"
	"time"
)

func Dial(address, password string) (*Conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	connection := newConnection(c, false)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := connection.SendCommand(ctx, command.Auth{Password: password})
	if err != nil {
		return nil, err
	}
	if !response.IsOk() {
		// Try to gracefully disconnect
		_, _ = connection.SendCommand(ctx, command.Exit{})
		return nil, errors.New("invalid authentication credentials")
	}

	return connection, nil
}
