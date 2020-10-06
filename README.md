# eslgo
eslgo is a [FreeSWITCH™](https://freeswitch.com/) ESL library for GoLang

## Install
```
go get github.com/percipia/eslgo
```
```
github.com/percipia/eslgo v1.2.1
```

## Overview
- Inbound ESL Connections
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
	"time"
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