package freeswitchesl

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"gitlab.percipia.com/libs/go/freeswitchesl/command/call"
	"log"
)

func (c *Conn) OriginateCall(ctx context.Context, aLeg, bLeg string, vars map[string]string) (string, error) {
	if vars == nil {
		vars = make(map[string]string)
	}
	vars["origination_uuid"] = uuid.New().String()

	_, err := c.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  fmt.Sprintf("%s%s %s", buildVars(vars), aLeg, bLeg),
		Background: true,
	})
	if err != nil {
		return vars["origination_uuid"], err
	}
	return vars["origination_uuid"], nil
}

func (c *Conn) HangupCall(ctx context.Context, uuid string) error {
	_, err := c.SendCommand(ctx, call.Hangup{
		UUID:  uuid,
		Cause: "NORMAL_CLEARING",
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

func (c *Conn) PlaybackAndWait(ctx context.Context, uuid, audioCommand string, times int) error {
	response, err := c.SendCommand(ctx, &call.Execute{
		UUID:    uuid,
		AppName: "playback",
		AppArgs: audioCommand,
		Loops:   times,
		Sync:    true,
	})
	if err != nil {
		return err
	}

	done := make(chan bool, 1)
	listenerID := c.RegisterEventListener(response.ChannelUUID(), func(event *Event) {
		log.Printf("%#v\n", event)
		if event.Headers.Get("Event-Name") == "PLAYBACK_STOP" {
			done <- true
		}
	})
	<-done
	c.RemoveEventListener(response.ChannelUUID(), listenerID)
	return nil
}
