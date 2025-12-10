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
	"strings"

	"github.com/percipia/eslgo/command"
	"github.com/percipia/eslgo/command/call"
)

// Leg This struct is used to specify the individual legs of a call for the originate helpers
type Leg struct {
	CallURL      string
	LegVariables map[string]string
}

// OriginateCall - Calls the originate function in FreeSWITCH. If you want variables for each leg independently set them in the aLeg and bLeg
// Arguments: ctx context.Context for supporting context cancellation, background bool should we wait for the origination to complete
// aLeg, bLeg Leg The aLeg and bLeg of the call respectively
// vars map[string]string, channel variables to be passed to originate for both legs, contained in {}
func (c *Conn) OriginateCall(ctx context.Context, background bool, aLeg, bLeg Leg, vars map[string]string) (*RawResponse, error) {
	if vars == nil {
		vars = make(map[string]string)
	}

	if _, ok := vars["origination_uuid"]; ok {
		// We cannot set origination uuid globally
		delete(vars, "origination_uuid")
	}

	response, err := c.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  fmt.Sprintf("%s%s %s", BuildVars("{%s}", vars), aLeg.String(), bLeg.String()),
		Background: background,
	})

	return response, err
}

// EnterpriseOriginateCall - Calls the originate function in FreeSWITCH using the enterprise method for calling multiple legs ":_:"
// If you want variables for each leg independently set them in the aLeg and bLeg strings
// Arguments: ctx context.Context for supporting context cancellation, background bool should we wait for the origination to complete
// vars map[string]string, channel variables to be passed to originate for both legs, contained in <>
// bLeg string The bLeg of the call
// aLegs ...string variadic argument for each aLeg to call
func (c *Conn) EnterpriseOriginateCall(ctx context.Context, background bool, vars map[string]string, bLeg Leg, aLegs ...Leg) (*RawResponse, error) {
	if len(aLegs) == 0 {
		return nil, errors.New("no aLeg specified")
	}

	if vars == nil {
		vars = make(map[string]string)
	}

	if _, ok := vars["origination_uuid"]; ok {
		// We cannot set origination uuid globally
		delete(vars, "origination_uuid")
	}

	var aLeg strings.Builder
	for i, leg := range aLegs {
		if i > 0 {
			aLeg.WriteString(":_:")
		}
		aLeg.WriteString(leg.String())
	}

	response, err := c.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  fmt.Sprintf("%s%s %s", BuildVars("<%s>", vars), aLeg.String(), bLeg.String()),
		Background: background,
	})

	return response, err
}

// BackgroundOriginateCall - Calls the originate function in FreeSWITCH asynchronously. If you want variables for each leg independently set them in the aLeg and bLeg
// Arguments: ctx context.Context for supporting context cancellation, background bool should we wait for the origination to complete
// aLeg, bLeg Leg The aLeg and bLeg of the call respectively
// vars map[string]string, channel variables to be passed to originate for both legs, contained in {}
func (c *Conn) BackgroundOriginateCall(ctx context.Context, background bool, aLeg, bLeg Leg, vars map[string]string) (*RawResponse, error) {
	if vars == nil {
		vars = make(map[string]string)
	}

	if _, ok := vars["origination_uuid"]; ok {
		// We cannot set origination uuid globally
		delete(vars, "origination_uuid")
	}

	response, err := c.SendCommand(ctx, command.API{
		Command:    "bgapi originate",
		Arguments:  fmt.Sprintf("%s%s %s", BuildVars("{%s}", vars), aLeg.String(), bLeg.String()),
		Background: background,
	})

	return response, err
}

// HangupCall - A helper to hangup a call asynchronously
func (c *Conn) HangupCall(ctx context.Context, uuid, cause string) error {
	_, err := c.SendCommand(ctx, call.Hangup{
		UUID:  uuid,
		Cause: cause,
		Sync:  false,
	})
	return err
}

// HangupCall - A helper to answer a call synchronously
func (c *Conn) AnswerCall(ctx context.Context, uuid string) error {
	_, err := c.SendCommand(ctx, &call.Execute{
		UUID:    uuid,
		AppName: "answer",
		Sync:    true,
	})
	return err
}

// String - Build the Leg string for passing to Bridge/Originate functions
func (l Leg) String() string {
	return fmt.Sprintf("%s%s", BuildVars("[%s]", l.LegVariables), l.CallURL)
}
