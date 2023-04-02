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

	prefix := bot.Prefix(".")
	robot.On(prefix.Cmd("ping").And(bot.AuthoredBy("gastrodon")), func(event *ifunny.ChatEvent) error {
		fmt.Printf("we got a ping from ourselves\n")
		return nil
	})

	us, _ := robot.Client.GetUser(compose.UserByNick("gastrodon"))
	robot.Subscribe(robot.Client.DMChannelName(us.ID))
	robot.Listen()
}
