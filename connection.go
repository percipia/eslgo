package freeswitchesl

import (
	"bufio"
	"context"
	"github.com/google/uuid"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
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
	eventListenerLock sync.RWMutex
	eventListeners    map[string]map[string]EventListener
	outbound          bool
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
	select {
	case response := <-c.responseChannels[TypeReply]:
		return response, nil
	case response := <-c.responseChannels[TypeAPIResponse]:
		return response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Conn) Close() {
	c.stopFunc()
	_ = c.conn.Close()
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
	channelUUID := event.Headers.Get("Unique-Id")
	if listeners, ok := c.eventListeners[channelUUID]; ok {
		for _, listener := range listeners {
			go listener(event)
		}
	}

	// Next call any listeners for a particular application
	appUUID := event.Headers.Get("Application-UUID")
	if listeners, ok := c.eventListeners[appUUID]; ok {
		for _, listener := range listeners {
			go listener(event)
		}
	}
}

func (c *Conn) eventLoop() {
	for {
		var event *Event
		var err error
		select {
		case raw := <-c.responseChannels[TypeEventPlain]:
			event, err = readPlainEvent(raw.Body)
		case raw := <-c.responseChannels[TypeEventXML]:
			event, err = readXMLEvent(raw.Body)
		case raw := <-c.responseChannels[TypeEventJSON]:
			event, err = readJSONEvent(raw.Body)
		case <-c.runningContext.Done():
			return
		}

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

		if responseChan, ok := c.responseChannels[response.Headers.Get("Content-Type")]; ok {
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
