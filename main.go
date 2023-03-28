package main

import (
	"fmt"
	"os"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny"
)

var bearer = ""
var userAgent = ""

func init() {
	bearer = os.Getenv("IFUNNY_BEARER")
	if bearer == "" {
		panic("IFUNNY_BEARER must be set")
	}

	userAgent = os.Getenv("IFUNNY_USER_AGENT")
	if userAgent == "" {
		panic("IFUNNY_USER_AGENT must be set")
	}
}

func main() {
	bot, err := bot.MakeBot(bearer, userAgent)
	if err != nil {
		panic(err)
	}

	contacts, err := bot.Chat.GetUsers(ifunny.Contacts(10_000))
	if err != nil {
		panic(err)
	}

	for _, contact := range contacts {
		bot.Subscribe(ifunny.DMChannelName(bot.Client.Self.ID, []string{contact.ID}))
	}

	bot.On(
		func(event ifunny.Event) bool { return event.Type() == 200 },
		func(event ifunny.Event) error {
			message := new(ifunny.ChatMessage)
			if err := event.Decode(message); err != nil {
				return nil
			}

			fmt.Printf("[%s %d] %s\n", message.User.Nick, message.Type, message.Text)
			return nil
		})

	bot.On(
		func(event ifunny.Event) bool { return event.Type() != 200 },
		func(event ifunny.Event) error {
			fmt.Printf("unknown event %d: %+v\n", event.Type(), event)
			return nil
		})

	for i := 0; i < 100; i++ {
		bot.Chat.HideChannel(fmt.Sprintf("foobar%d", i))
	}

	bot.Listen()
}
