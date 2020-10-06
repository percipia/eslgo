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

type Auth struct {
	User     string
	Password string
}

func (auth Auth) BuildMessage() string {
	if len(auth.User) > 0 {
		return fmt.Sprintf("userauth %s:%s", auth.User, auth.Password)
	}
	return fmt.Sprintf("auth %s", auth.Password)
}
