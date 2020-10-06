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
package eslgo

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_BuildVars(t *testing.T) {
	vars := BuildVars("{%s}", map[string]string{
		"origination_caller_name":   "test",
		"origination_caller_number": "1234",
		"origination_callee_name":   "John Doe",
		"origination_callee_number": "7100",
	})

	// Contains since order is not guaranteed when iterating over maps
	assert.Contains(t, vars, "origination_caller_name=test")
	assert.Contains(t, vars, "origination_caller_number=1234")
	assert.Contains(t, vars, "origination_callee_name='John Doe'")
	assert.Contains(t, vars, "origination_callee_number=7100")

	// Ensure the formatting elements are contained in the string
	assert.Equal(t, 3, strings.Count(vars, ","))
	assert.True(t, strings.HasPrefix(vars, "{"))
	assert.True(t, strings.HasSuffix(vars, "}"))
}
