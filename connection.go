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
	"bufio"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/percipia/eslgo/command"
	"log"
	"net"
	"net/textproto"
	"sync"
	"time"
)

type Conn struct {
	conn              net.Conn
	reader            *bufio.Reader
	header            *textproto.Reader
	writeLock         sync.Mutex
	runningContext    context.Context
	stopFunc          func()
	responseChannels  map[string]chan *RawResponse
	responseChanMutex sync.RWMutex
	eventListenerLock sync.RWMutex
	eventListeners    map[string]map[string]EventListener
	outbound          bool
	closeOnce         sync.Once
}

const EndOfMessage = "\r\n\r\n"

func newConnection(c net.Conn, outbound bool) *Conn {
	reader := bufio.NewReader(c)
	header := textproto.NewReader(reader)

	runningContext, stop := context.WithCancel(context.Background())

	instance := &Conn{
		conn:   c,
		reader: reader,
		header: header,
		responseChannels: map[string]chan *RawResponse{
			TypeReply:       make(chan *RawResponse),
			TypeAPIResponse: make(chan *RawResponse),
			TypeEventPlain:  make(chan *RawResponse),
			TypeEventXML:    make(chan *RawResponse),
			TypeEventJSON:   make(chan *RawResponse),
			TypeAuthRequest: make(chan *RawResponse),
			TypeDisconnect:  make(chan *RawResponse),
		},
		runningContext: runningContext,
		stopFunc:       stop,
		eventListeners: make(map[string]map[string]EventListener),
		outbound:       outbound,
	}
	go instance.receiveLoop()
	go instance.eventLoop()
	return instance
}

func (c *Conn) RegisterEventListener(channelUUID string, listener EventListener) string {
	c.eventListenerLock.Lock()
	defer c.eventListenerLock.Unlock()

	id := uuid.New().String()
	if _, ok := c.eventListeners[channelUUID]; ok {
		c.eventListeners[channelUUID][id] = listener
	} else {
		c.eventListeners[channelUUID] = map[string]EventListener{id: listener}
	}
	return id
}

func (c *Conn) RemoveEventListener(channelUUID string, id string) {
	c.eventListenerLock.Lock()
	defer c.eventListenerLock.Unlock()

	if listeners, ok := c.eventListeners[channelUUID]; ok {
		delete(listeners, id)
	}
}

func (c *Conn) SendCommand(ctx context.Context, command command.Command) (*RawResponse, error) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	if deadline, ok := ctx.Deadline(); ok {
		_ = c.conn.SetWriteDeadline(deadline)
	}
	_, err := c.conn.Write([]byte(command.BuildMessage() + EndOfMessage))
	if err != nil {
		return nil, err
	}

	// Get response
	c.responseChanMutex.RLock()
	defer c.responseChanMutex.RUnlock()
	select {
	case response := <-c.responseChannels[TypeReply]:
		if response == nil {
			// We only get nil here if the channel is closed
			return nil, errors.New("connection closed")
		}
		return response, nil
	case response := <-c.responseChannels[TypeAPIResponse]:
		if response == nil {
			// We only get nil here if the channel is closed
			return nil, errors.New("connection closed")
		}
		return response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Conn) Close() {
	c.closeOnce.Do(c.close)
}

func (c *Conn) close() {
	c.stopFunc()
	_ = c.conn.Close()
	c.responseChanMutex.Lock()
	defer c.responseChanMutex.Unlock()
	for key, responseChan := range c.responseChannels {
		close(responseChan)
		delete(c.responseChannels, key)
	}
}

func (c *Conn) callEventListener(event *Event) {
	c.eventListenerLock.RLock()
	defer c.eventListenerLock.RUnlock()

	// First check if there are any general event listener
	if listeners, ok := c.eventListeners[EventListenAll]; ok {
		for _, listener := range listeners {
			go listener(event)
		}
	}

	// Next call any listeners for a particular channel
	if event.HasHeader("Unique-Id") {
		channelUUID := event.GetHeader("Unique-Id")
		if listeners, ok := c.eventListeners[channelUUID]; ok {
			for _, listener := range listeners {
				go listener(event)
			}
		}
	}

	// Next call any listeners for a particular application
	if event.HasHeader("Application-UUID") {
		appUUID := event.GetHeader("Application-UUID")
		if listeners, ok := c.eventListeners[appUUID]; ok {
			for _, listener := range listeners {
				go listener(event)
			}
		}
	}

	// Next call any listeners for a particular job
	if event.HasHeader("Job-UUID") {
		jobUUID := event.GetHeader("Job-UUID")
		if listeners, ok := c.eventListeners[jobUUID]; ok {
			for _, listener := range listeners {
				go listener(event)
			}
		}
	}
}

func (c *Conn) eventLoop() {
	for {
		var event *Event
		var err error
		c.responseChanMutex.RLock()
		select {
		case raw := <-c.responseChannels[TypeEventPlain]:
			if raw == nil {
				// We only get nil here if the channel is closed
				c.responseChanMutex.RUnlock()
				return
			}
			event, err = readPlainEvent(raw.Body)
		case raw := <-c.responseChannels[TypeEventXML]:
			if raw == nil {
				// We only get nil here if the channel is closed
				c.responseChanMutex.RUnlock()
				return
			}
			event, err = readXMLEvent(raw.Body)
		case raw := <-c.responseChannels[TypeEventJSON]:
			if raw == nil {
				// We only get nil here if the channel is closed
				c.responseChanMutex.RUnlock()
				return
			}
			event, err = readJSONEvent(raw.Body)
		case <-c.runningContext.Done():
			c.responseChanMutex.RUnlock()
			return
		}
		c.responseChanMutex.RUnlock()

		if err != nil {
			log.Printf("Error parsing event\n%s\n", err.Error())
			continue
		}

		c.callEventListener(event)
	}
}

func (c *Conn) receiveLoop() {
	for {
		response, err := c.readResponse()
		if err != nil {
			break
		}

		c.responseChanMutex.RLock()
		responseChan, ok := c.responseChannels[response.GetHeader("Content-Type")]
		if !ok && len(c.responseChannels) <= 0 {
			// We must have shutdown!
			break
		}
		c.responseChanMutex.RUnlock()
		if ok {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			select {
			case responseChan <- response:
			case <-c.runningContext.Done():
				cancel()
				return
			case <-ctx.Done():
				log.Printf("No one to handle response %v\n", response)
			}
			cancel()
		}
	}
}
