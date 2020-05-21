package freeswitchesl

import (
	"bufio"
	"context"
	"errors"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"log"
	"net"
	"net/textproto"
	"sync"
	"time"
)

type Conn struct{
	conn             net.Conn
	reader           *bufio.Reader
	header           *textproto.Reader
	writeLock		 sync.Mutex
	runningContext   context.Context
	stopFunc         func()
	responseChannels map[string]chan *RawResponse
}

const EndOfMessage = "\r\n\r\n"

func newConnection(c net.Conn) *Conn {
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
		stopFunc: stop,
	}
	go instance.receiveLoop()
	return instance
}

func (c *Conn) Close() {
	c.stopFunc()
	_ = c.conn.Close()
}

func (c *Conn) sendCommand(ctx context.Context, command command.Command) (*RawResponse, error) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	if deadline, ok := ctx.Deadline(); ok {
		_ = c.conn.SetWriteDeadline(deadline)
	}
	_, err := c.conn.Write([]byte(command.String() + EndOfMessage))
	if err != nil {
		return nil, err
	}
	return c.waitFor(ctx, TypeReply)
}

func (c *Conn) waitFor(ctx context.Context, responseType string) (*RawResponse, error) {
	if responseChan, ok := c.responseChannels[responseType]; ok {
		select {
		case response := <- responseChan:
			return response, nil
		case <- ctx.Done():
			return nil, ctx.Err()
		}
	}
	return nil, errors.New("no such response content type")
}

func (c *Conn) eventLoop() {
	for {
		var raw *RawResponse
		select {
		case raw = <- c.responseChannels[TypeEventPlain]:
		case raw = <- c.responseChannels[TypeEventXML]:
		case raw = <- c.responseChannels[TypeEventJSON]:
		case <- c.runningContext.Done():
			return
		}

		log.Printf("Event %s\n", string(raw.Body))
	}
}

func (c *Conn) receiveLoop() {
	for {
		response, err := c.readResponse()
		if err != nil {
			break
		}

		if responseChan, ok := c.responseChannels[response.Headers.Get("Content-Type")]; ok {
			ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
			select {
			case responseChan <- response:
			case <- c.runningContext.Done():
				return
			case <- ctx.Done():
				log.Printf("No one to handle response %#v\n", response)
			}
			cancel()
		}
	}
}