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

type API struct {
	Command    string
	Arguments  string
	Background bool
}

func (api API) BuildMessage() string {
	if api.Background {
		return fmt.Sprintf("bgapi %s %s", api.Command, api.Arguments)
	}
	return fmt.Sprintf("api %s %s", api.Command, api.Arguments)
}
