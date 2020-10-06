/*
 * Copyright (c) 2020 Percipia
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Contributor(s):
 * Andrew Querol <aquerol@percipia.com>
 */
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
