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

type Filter struct {
	Delete      bool
	EventHeader string
	FilterValue string
}

func (f Filter) BuildMessage() string {
	if f.Delete {
		if len(f.FilterValue) > 0 {
			// Clear just the specific header value
			return fmt.Sprintf("filter delete %s %s", f.EventHeader, f.FilterValue)
		}
		// Clears all filters for the header
		return fmt.Sprintf("filter delete %s", f.EventHeader)
	}
	return fmt.Sprintf("filter %s %s", f.EventHeader, f.FilterValue)
}
