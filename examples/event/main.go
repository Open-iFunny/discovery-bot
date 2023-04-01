package main

import (
	"fmt"
	"os"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny"
)

var bearer = os.Getenv("IFUNNY_BEARER")
var userAgent = os.Getenv("IFUNNY_USER_AGENT")

func main() {
	bot, err := bot.MakeBot(bearer, userAgent)
	if err != nil {
		panic(err)
	}

	bot.Subscribe("apitools")
	bot.OnMessage(func(event *ifunny.ChatEvent) error {
		fmt.Printf("[%s] %s\n", event.User.Nick, event.Text)
		return nil
	})

	bot.Listen()
}
