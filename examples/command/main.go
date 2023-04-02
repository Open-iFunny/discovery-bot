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
	robot, err := bot.MakeBot(bearer, userAgent)
	if err != nil {
		panic(err)
	}

	prefix := bot.Prefix(".")
	robot.On(prefix.Cmd("ping").And(bot.AuthoredBy("gastrodon")), func(event *ifunny.ChatEvent) error {
		fmt.Printf("we got a ping from ourselves\n")
		return nil
	})

	robot.On(prefix.Cmd("ping").Not(bot.AuthoredBy("gastrodon")), func(event *ifunny.ChatEvent) error {
		fmt.Printf("we got a ping from somebody else\n")
		return nil
	})

	robot.Chat.OnChannelUpdate(func(eventType int, channel *ifunny.ChatChannel) error {
		if eventType == ifunny.EVENT_JOIN {
			robot.Subscribe(channel.Name)
		}

		return nil
	})

	robot.Listen()
}
