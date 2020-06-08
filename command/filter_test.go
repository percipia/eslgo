package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter_BuildMessage(t *testing.T) {
	assert.Equal(t, "filter variable_domain_name 192.168.1.1", Filter{
		EventHeader: "variable_domain_name",
		FilterValue: "192.168.1.1",
	}.BuildMessage())
	assert.Equal(t, "filter delete variable_domain_name 192.168.1.1", Filter{
		Delete:      true,
		EventHeader: "variable_domain_name",
		FilterValue: "192.168.1.1",
	}.BuildMessage())
}
