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
	"fmt"
	"time"
)

type Linger struct {
	Enabled bool
	Seconds time.Duration
}

func (l Linger) BuildMessage() string {
	if l.Enabled {
		if l.Seconds > 0 {
			return fmt.Sprintf("linger %d", l.Seconds)
		}
		return "linger"
	}
	return "nolinger"
}
