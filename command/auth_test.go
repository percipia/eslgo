package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestAuthMessage     = `auth testing123`
	TestUserAuthMessage = `userauth testuser:testing123`
)

func TestAuth_BuildMessage(t *testing.T) {
	auth := Auth{
		Password: "testing123",
	}
	assert.Equal(t, TestAuthMessage, auth.BuildMessage())
}

func TestAuth_BuildMessage_User(t *testing.T) {
	auth := Auth{
		User:     "testuser",
		Password: "testing123",
	}
	assert.Equal(t, TestUserAuthMessage, auth.BuildMessage())
}
