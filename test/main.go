package main

import (
	"fmt"
	"gitlab.percipia.com/libs/go/freeswitchesl"
	"log"
)

func main() {
	log.Fatalln(freeswitchesl.ListenAndServe(":8084", handleConnection))
}

func handleConnection(conn *freeswitchesl.Conn, response *freeswitchesl.RawResponse) {
	fmt.Printf("Got connection! %#v\n%#v\n\n", response.Headers, string(response.Body))
}