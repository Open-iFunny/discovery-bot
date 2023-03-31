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
	client, _ := ifunny.MakeClient(bearer, userAgent)
	chat, _ := client.Chat()

	channels, err := client.GetChannels(compose.ChatsTrending)
	if err != nil {
		panic(err)
	}

	fmt.Printf("got %d trendy chat channels!\n", len(channels))

	messages, _, _, err := chat.ListMessages(compose.ListMessages("apitools", 10, compose.NoPage()))
	if err != nil {
		panic(err)
	}

	fmt.Printf("got %d messages from apitools!\n", len(messages))
}
