package main

import (
	"fmt"
	"os"

	"github.com/gastrodon/popplio/ifunny"
	"github.com/gastrodon/popplio/ifunny/compose"
)

var bearer = os.Getenv("IFUNNY_BEARER")
var userAgent = os.Getenv("IFUNNY_USER_AGENT")

func main() {
	client, err := ifunny.MakeClient(bearer, userAgent)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s [id: %s] authed!\n", client.Self.Nick, client.Self.ID)

	chat, err := client.Chat()
	if err != nil {
		panic(err)
	}

	chat.Subscribe(compose.EventsIn("apitools"), func(eventType int, kwargs map[string]interface{}) error {
		fmt.Printf("RECV event of type %d: %+v\n", eventType, kwargs)
		return nil
	})

	<-make(chan bool)
}
