package freeswitchesl

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"gitlab.percipia.com/libs/go/freeswitchesl/command/call"
	"io"
	"log"
	"strings"
)

func (c *Conn) EnableEvents(ctx context.Context) error {
	var err error
	if c.outbound {
		_, err = c.SendCommand(ctx, command.MyEvents{
			Format: "plain",
		})
	} else {
		_, err = c.SendCommand(ctx, command.Event{
			Format: "plain",
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

func (c *Conn) OriginateCall(ctx context.Context, aLeg, bLeg string, vars map[string]string) (string, *RawResponse, error) {
	if vars == nil {
		vars = make(map[string]string)
	}
	if _, ok := vars["origination_uuid"]; !ok {
		vars["origination_uuid"] = uuid.New().String()
	}
	response, err := c.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  fmt.Sprintf("%s%s %s", buildVars("{%s}", vars), aLeg, bLeg),
		Background: true,
	})
	if err != nil {
		return vars["origination_uuid"], err
	}
	return vars["origination_uuid"], nil
}

func (c *Conn) EnterpriseOriginateCall(ctx context.Context, vars map[string]string, bLeg string, aLegs ...string) (string, error) {
	if vars == nil {
		vars = make(map[string]string)
	}
	vars["origination_uuid"] = uuid.New().String()

	if len(aLegs) == 0 {
		return "", errors.New("no aLeg specified")
	}
	aLeg := strings.Join(aLegs, ":_:")

	_, err := c.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  fmt.Sprintf("%s%s %s", buildVars("<%s>", vars), aLeg, bLeg),
		Background: true,
	})
	if err != nil {
		return vars["origination_uuid"], response, err
	}
	return vars["origination_uuid"], response, nil
}

func (c *Conn) HangupCall(ctx context.Context, uuid, cause string) error {
	_, err := c.SendCommand(ctx, call.Hangup{
		UUID:  uuid,
		Cause: cause,
		Sync:  false,
	})
	return err
}

func (c *Conn) AnswerCall(ctx context.Context, uuid string) error {
	_, err := c.SendCommand(ctx, &call.Execute{
		UUID:    uuid,
		AppName: "answer",
		Sync:    true,
	})
	return err
}

// Playback - Executes the mod_dptools playback app
func (c *Conn) Playback(ctx context.Context, uuid, audioArgs string, times int, wait bool) error {
	return c.audioCommand(ctx, "playback", uuid, audioArgs, times, wait)
}

// Say - Executes the mod_dptools say app
func (c *Conn) Say(ctx context.Context, uuid, audioArgs string, times int, wait bool) error {
	return c.audioCommand(ctx, "say", uuid, audioArgs, times, wait)
}

// Speak - Executes the mod_dptools speak app
func (c *Conn) Speak(ctx context.Context, uuid, audioArgs string, times int, wait bool) error {
	return c.audioCommand(ctx, "speak", uuid, audioArgs, times, wait)
}

// WaitForDTMF, waits for a DTMF event. Requires events to be enabled!
func (c *Conn) WaitForDTMF(ctx context.Context, uuid string) (byte, error) {
	done := make(chan byte, 1)
	listenerID := c.RegisterEventListener(uuid, func(event *Event) {
		if event.GetName() == "DTMF" {
			dtmf := event.GetHeader("DTMF-Digit")
			if len(dtmf) > 0 {
				done <- dtmf[0]
			}
			done <- 0
		}
	})
	defer c.RemoveEventListener(uuid, listenerID)
	defer close(done)

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
func (c *Conn) audioCommand(ctx context.Context, command, uuid, audioArgs string, times int, wait bool) error {
	response, err := c.SendCommand(ctx, &call.Execute{
		UUID:    uuid,
		AppName: command,
		AppArgs: audioArgs,
		Loops:   times,
		Sync:    wait,
	})
	if err != nil {
		return err
	}
	if !response.IsOk() {
		return errors.New(command + " response is not okay")
	}
	return nil
}
