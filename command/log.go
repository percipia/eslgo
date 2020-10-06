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

import "fmt"

type Log struct {
	Enabled bool
	Level   int
}

func (l Log) BuildMessage() string {
	if l.Enabled {
		return fmt.Sprintf("log %d", l.Level)
	}
	return "nolog"
}
