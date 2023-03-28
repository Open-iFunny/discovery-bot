package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gastrodon/popplio/ifunny"
)

func main() {
	bearer := os.Getenv("IFUNNY_BEARER")
	if bearer == "" {
		panic("IFUNNY_BEARER must be set")
	}

	userAgent := os.Getenv("IFUNNY_USER_AGENT")
	if userAgent == "" {
		panic("IFUNNY_USER_AGENT must be set")
	}

	client, _ := ifunny.MakeClient(bearer, userAgent)
	client.User(ifunny.UserAccount)

	<-time.After(5 * time.Second)

	chat, _ := client.Chat()

	go func() {
		channel, _, _ := chat.GetChannel(client.ChannelDM("5396c1348ea6b8dc5a8b456c"))

		iterMessage, _ := chat.IterMessage(ifunny.MessageIn(channel.Name))
		for message := range iterMessage {
			fmt.Printf("recv: %+v\n", message)
		}
	}()

	go func() {
		channel, _, _ := chat.GetChannel(ifunny.ChannelName("yoloswaggin"))

		iterMessage, ubsubscribe := chat.IterMessage(ifunny.MessageIn(channel.Name))
		go func() {
			<-time.After(5 * time.Second)
			fmt.Println("unsubscribing yoloswaggin")
			ubsubscribe()
		}()

		for message := range iterMessage {
			fmt.Printf("recv: %+v\n", message)
		}
	}()

	<-time.After(40 * time.Second)
}
