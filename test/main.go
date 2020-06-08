package main

import (
	"context"
	"fmt"
	"gitlab.percipia.com/libs/go/freeswitchesl"
	"gitlab.percipia.com/libs/go/freeswitchesl/command"
	"log"
	"time"
)

func main() {
	log.Fatalln(freeswitchesl.ListenAndServe(":8084", handleConnection))
}

func handleConnection(conn *freeswitchesl.Conn, response *freeswitchesl.RawResponse) {
	fmt.Printf("Got connection! %#v\n", response)
	conn.SendCommand(context.Background(), command.Event{
		Format: "plain",
		Listen: []string{"ALL"},
	})
	conn.SendCommand(context.Background(), command.API{
		Command:   "originate",
		Arguments: "user/100 &park()",
	})
	time.Sleep(60 * time.Second)
}
