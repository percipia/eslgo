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
	"io"
	"net/textproto"
	"strings"
)

/*
 * FreeSWITCH allows characters in variable names that net/textproto rejects as
 * malformed MIME header names, for example "variable_fish,(sticks)~_1". The Go
 * standard library validates header names against RFC 7230 token rules, so a
 * single such variable makes ReadMIMEHeader fail and the whole event is
 * discarded.
 *
 * readESLMIMEHeader is a lenient replacement that follows the same wire format
 * (name, colon, optional space, value, terminated by a blank line) but accepts
 * any non-colon, non-space byte in the header name. Names are still canonicalized
 * with textproto.CanonicalMIMEHeaderKey when they are valid tokens, so existing
 * code using Get/HasHeader keeps working unchanged.
 */
func readESLMIMEHeader(r *textproto.Reader) (textproto.MIMEHeader, error) {
	header := make(textproto.MIMEHeader)
	for {
		line, err := r.ReadLine()
		if err != nil {
			// A truncated header block is still worth returning: the caller can use
			// whatever was parsed instead of dropping the entire event.
			if err == io.EOF && len(header) > 0 {
				return header, nil
			}
			return header, err
		}
		// Blank line terminates the header block.
		if len(line) == 0 {
			return header, nil
		}

		key, value, ok := splitHeaderLine(line)
		if !ok {
			// No colon at all, there is nothing sane to do with this line. Skip it
			// rather than discarding every header that followed it.
			continue
		}
		header[canonicalESLHeaderKey(key)] = append(header[canonicalESLHeaderKey(key)], value)
	}
}

// splitHeaderLine splits "Name: value" on the first colon. Leading whitespace in
// the value is trimmed, matching net/textproto behavior.
func splitHeaderLine(line string) (key, value string, ok bool) {
	idx := strings.IndexByte(line, ':')
	if idx <= 0 {
		return "", "", false
	}
	key = line[:idx]
	value = strings.TrimLeft(line[idx+1:], " \t")
	return key, value, true
}

/*
 * canonicalESLHeaderKey canonicalizes a header name the same way net/textproto
 * does when the name is a valid HTTP token, and returns it untouched otherwise.
 * This keeps standard FreeSWITCH headers such as "Event-Name" and "Unique-ID"
 * addressable through the existing Get/HasHeader helpers, while preserving the
 * exact spelling of variable names that contain characters like , ( ) ~ so
 * callers can still look them up by their real FreeSWITCH name.
 */
func canonicalESLHeaderKey(key string) string {
	if isValidHeaderToken(key) {
		return textproto.CanonicalMIMEHeaderKey(key)
	}
	return key
}

// isValidHeaderToken reports whether every byte of key is a valid RFC 7230 token
// character, which is the same set net/textproto accepts.
func isValidHeaderToken(key string) bool {
	if len(key) == 0 {
		return false
	}
	for i := 0; i < len(key); i++ {
		if !isTokenChar(key[i]) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c >= '0' && c <= '9':
		return true
	}
	switch c {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
		return true
	}
	return false
}
