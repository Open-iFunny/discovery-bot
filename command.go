package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny/compose"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

var okID = os.Getenv("IFUNNY_ADMIN_ID")

func doHistSnapshot(channels chan<- string) func(ctx bot.Context) error {
	return func(ctx bot.Context) error {
		event, err := ctx.Event()
		if err != nil {
			return err
		}

		parts := strings.Split(event.Text, " ")
		if len(parts) < 2 {
			return ctx.Send(".snap <channel-id>")
		}

		channel, err := ctx.Robot().Chat.GetChannel(compose.GetChannel(parts[1]))
		if err != nil {
			ctx.Send(err.Error())
			return err
		}

		go func() {
			channels <- channel.Name
			ctx.Send("snap" + channel.Name + " begin")
		}()

		return ctx.Send("snap " + channel.Name + " enqueue'd")
	}
}

func onCommand(robot *bot.Bot) error {
	byID := bot.AuthoredBy(okID)
	prefix := bot.Prefix(".")

	startTime := time.Now()
	robot.On(prefix.Cmd("uptime").And(byID), func(ctx bot.Context) error {
		event, err := ctx.Event()
		if err != nil {
			return err
		}

		robot.Log.WithFields(logrus.Fields{"cmd": "uptime", "caller": event.User.Nick}).Trace("call")
		return ctx.Send(fmt.Sprintf("Up for %d hours", int(math.Floor(time.Since(startTime).Hours()))))
	})

	robot.On(prefix.Cmd("snap").And(byID), doHistSnapshot(histChan))

	return nil
}
