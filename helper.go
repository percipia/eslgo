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
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/percipia/eslgo/command"
	"github.com/percipia/eslgo/command/call"
)

func (c *Conn) EnableEvents(ctx context.Context, format ...string) error {
	var err error
	eventFormat := "plain" // default to plain text
	if len(format) > 0 && format[0] != "" {
		eventFormat = format[0]
	}

	if c.outbound {
		_, err = c.SendCommand(ctx, command.MyEvents{
			Format: eventFormat,
		})
	} else {
		_, err = c.SendCommand(ctx, command.Event{
			Format: eventFormat,
			Listen: []string{"all"},
		})
	}
	return err
}

// DebugEvents - A helper that will output all events to a logger
func (c *Conn) DebugEvents(w io.Writer) string {
	logger := log.New(w, "EventLog: ", log.LstdFlags|log.Lmsgprefix)
	return c.RegisterEventListener(EventListenAll, func(event *Event) {
		logger.Println(event)
	})
}

func (c *Conn) DebugOff(id string) {
	c.RemoveEventListener(EventListenAll, id)
}

// Phrase - Executes the mod_dptools phrase app
func (c *Conn) Phrase(ctx context.Context, uuid, macro string, times int, wait bool) (*RawResponse, error) {
	return c.audioCommand(ctx, "phrase", uuid, macro, times, wait)
}

// PhraseWithArg - Executes the mod_dptools phrase app with arguments
func (c *Conn) PhraseWithArg(ctx context.Context, uuid, macro string, argument interface{}, times int, wait bool) (*RawResponse, error) {
	return c.audioCommand(ctx, "phrase", uuid, fmt.Sprintf("%s,%v", macro, argument), times, wait)
}

// Playback - Executes the mod_dptools playback app
func (c *Conn) Playback(ctx context.Context, uuid, audioArgs string, times int, wait bool) (*RawResponse, error) {
	return c.audioCommand(ctx, "playback", uuid, audioArgs, times, wait)
}

// Say - Executes the mod_dptools say app
func (c *Conn) Say(ctx context.Context, uuid, audioArgs string, times int, wait bool) (*RawResponse, error) {
	return c.audioCommand(ctx, "say", uuid, audioArgs, times, wait)
}

// Speak - Executes the mod_dptools speak app
func (c *Conn) Speak(ctx context.Context, uuid, audioArgs string, times int, wait bool) (*RawResponse, error) {
	return c.audioCommand(ctx, "speak", uuid, audioArgs, times, wait)
}

// WaitForDTMF, waits for a DTMF event. Requires events to be enabled!
func (c *Conn) WaitForDTMF(ctx context.Context, uuid string) (byte, error) {
	done := make(chan byte, 1)
	listenerID := c.RegisterEventListener(uuid, func(event *Event) {
		if event.GetName() == "DTMF" {
			dtmf := event.GetHeader("DTMF-Digit")
			if len(dtmf) > 0 {
				select {
				case done <- dtmf[0]:
				default:
				}
			} else {
				select {
				case done <- 0:
				default:
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	})
	defer func() {
		c.RemoveEventListener(uuid, listenerID)
		close(done)
	}()

	select {
	case digit := <-done:
		if digit != 0 {
			return digit, nil
		}
		return digit, errors.New("invalid DTMF digit received")
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// Helper for mod_dptools apps since they are very similar in invocation
func (c *Conn) audioCommand(ctx context.Context, command, uuid, audioArgs string, times int, wait bool) (*RawResponse, error) {
	response, err := c.SendCommand(ctx, &call.Execute{
		UUID:    uuid,
		AppName: command,
		AppArgs: audioArgs,
		Loops:   times,
		Sync:    wait,
	})
	if err != nil {
		return response, err
	}
	if !response.IsOk() {
		return response, errors.New(command + " response is not okay")
	}
	return response, nil
}
