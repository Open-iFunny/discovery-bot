package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny/compose"

	_ "github.com/go-sql-driver/mysql"
)

var okID = os.Getenv("IFUNNY_ADMIN_ID")

func cmdHistSnapshot(channels chan<- string) func(ctx bot.Context) error {
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

func cmdUptime(start time.Time) func(bot.Context) error {
	return func(ctx bot.Context) error {
		return ctx.Send(fmt.Sprintf("Up for %d hours", int(math.Floor(time.Since(start).Hours()))))
	}
}

func onCommand(_ *sql.DB, robot *bot.Bot) error {
	byID := bot.AuthoredBy(okID)
	prefix := bot.Prefix(".")
	robot.On(prefix.Cmd("uptime").And(byID), cmdUptime(time.Now()))
	robot.On(prefix.Cmd("snap").And(byID), cmdHistSnapshot(collectChannel))

	return nil
}
