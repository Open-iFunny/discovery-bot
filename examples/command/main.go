package main

import (
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

	prefix := bot.Prefix(".")
	robot.On(prefix.Cmd("ping").And(bot.AuthoredBy("gastrodon")), func(ctx bot.Context) error {
		if channel, err := ctx.Channel(); err != nil {
			return err
		} else {
			return robot.Chat.Publish(compose.MessageTo(channel.Name, "Hi:)"))
		}
	})

	robot.On(prefix.Cmd("ping").Not(bot.AuthoredBy("gastrodon")), func(ctx bot.Context) error {
		if channel, err := ctx.Channel(); err != nil {
			return err
		} else {
			return robot.Chat.Publish(compose.MessageTo(channel.Name, "I don't trust like that..."))
		}
	})

	robot.Chat.OnChannelUpdate(func(eventType int, channel *ifunny.ChatChannel) error {
		if eventType == ifunny.EVENT_JOIN {
			robot.Subscribe(channel.Name)
		}

		return nil
	})

	robot.Listen()
}
