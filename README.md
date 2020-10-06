# eslgo
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/percipia/eslgo/Go)
[![Total alerts](https://img.shields.io/lgtm/alerts/g/percipia/eslgo.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/percipia/eslgo/alerts/)
[![GitHub license](https://img.shields.io/github/license/percipia/eslgo)](https://github.com/percipia/eslgo/blob/v1/LICENSE)

eslgo is a [FreeSWITCH™](https://freeswitch.com/) ESL library for GoLang.
eslgo was written from the ground up in idiomatic Go fo use in our production products tested handling thousands of calls per second.

## Install
```
go get github.com/percipia/eslgo
```
```
github.com/percipia/eslgo v1.2.1
```

## Overview
- Inbound ESL Connection
- Outbound ESL Server
- Context Support
- Basic Helpers for common tasks
  - DTMF
  - Call origination
  - Call answer/hangup
  - Audio playback

## Examples
There are some buildable examples under the `example` directory as well
### Outbound ESL Server
```go
package main

import (
	"context"
	"fmt"
	"github.com/percipia/eslgo"
	"log"
)

func main() {
	log.Fatalln(eslgo.ListenAndServe(":8084", handleConnection))
}

func handleConnection(ctx context.Context, conn *eslgo.Conn, response *eslgo.RawResponse) {
	fmt.Printf("Got connection! %#v\n", response)
	_ = conn.EnableEvents(ctx)
	originationUUID, response, err := conn.OriginateCall(ctx, "user/100", "&playback(misc/ivr-to_hear_screaming_monkeys.wav)", map[string]string{})
	fmt.Println("Call Originated: ", originationUUID, response, err)
}
```
## Inbound ESL Client
```go
package main

import (
	"context"
	"fmt"
	"github.com/percipia/eslgo"
	"time"
)

func main() {
	conn, err := eslgo.Dial("127.0.0.1", "ClueCon", func() {
		fmt.Println("Inbound Connection Done")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
	defer cancel()

	_ = conn.EnableEvents(ctx)
	originationUUID, response, err := conn.OriginateCall(ctx, "user/100", "&playback(misc/ivr-to_hear_screaming_monkeys.wav)", map[string]string{})
	fmt.Println("Call Originated: ", originationUUID, response, err)

	time.Sleep(60 * time.Second)
}
```