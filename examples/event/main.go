package main

import (
	"fmt"
	"os"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny"
	"github.com/gastrodon/popplio/ifunny/compose"
)

var bearer = os.Getenv("IFUNNY_BEARER")
var userAgent = os.Getenv("IFUNNY_USER_AGENT")

func main() {
	robot, err := bot.MakeBot(bearer, userAgent)
	if err != nil {
		panic(err)
	}

	robot.OnMessage(func(ctx bot.Context) error {
		channel, err := ctx.Channel()
		if err != nil {
			return err
		}

		caller, err := ctx.Caller()
		if err != nil {
			return err
		}

		event, err := ctx.Event()
		if err != nil {
			return err
		}

		switch channel.Type {
		case compose.ChannelPublic, compose.ChannelPrivate:
			fmt.Printf("[%s | %s] %s\n", channel.Title, caller.Nick, event.Text)
		case compose.ChannelDM:
			fmt.Printf("[%s] %s\n", caller.Nick, event.Text)
		}

		return nil
	})

	robot.Chat.OnChannelUpdate(func(eventType int, channel *ifunny.ChatChannel) error {
		robot.Subscribe(channel.Name)
		return nil
	})

	robot.Listen()
}
