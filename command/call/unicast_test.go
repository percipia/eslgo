package call

import (
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
)

var (
	UnicastMessage = strings.ReplaceAll(`sendmsg none
Call-Command: unicast
Flags: native
Local-Ip: 192.168.1.100
Local-Port: 8025
Remote-Ip: 192.168.1.101
Remote-Port: 8026
Transport: tcp`, "\n", "\r\n")
)

func TestUnicast_BuildMessage(t *testing.T) {
	testLocal, _ := net.ResolveTCPAddr("tcp", "192.168.1.100:8025")
	testRemote, _ := net.ResolveTCPAddr("tcp", "192.168.1.101:8026")
	unicast := Unicast{
		UUID:   "none",
		Local:  testLocal,
		Remote: testRemote,
		Flags:  "native",
	}
	assert.Equal(t, UnicastMessage, unicast.BuildMessage())
}
