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
			return ctx.Send(".snap <channel-id>\nEnqueue snapshotting a chat channel")
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

func cmdInvite(ctx bot.Context) error {
	event, err := ctx.Event()
	if err != nil {
		return err
	}

	parts := strings.Split(event.Text, " ")
	switch len(parts) {
	case 2:
		caller, err := ctx.Caller()
		if err != nil {
			return err
		}

		if err := ctx.Robot().Chat.Call(compose.Invite(parts[1], []string{caller.ID}), nil); err != nil {
			return err
		}

		return ctx.Send(fmt.Sprintf("invited you to %s!", parts[1]))
	case 3:
		target, err := ctx.Robot().Client.GetUser(compose.UserByNick(parts[2]))
		if err != nil {
			return err
		}

		if err := ctx.Robot().Chat.Call(compose.Invite(parts[1], []string{target.ID}), nil); err != nil {
			return err
		}

		return ctx.Send(fmt.Sprintf("invited %s to %s!", target.Nick, parts[1]))
	default:
		return ctx.Send(".invite <channel-id> [user-nick]\nInvite yourself or another user to a chat channel")
	}
}

func onCommand(_ *sql.DB, robot *bot.Bot) error {
	byID := bot.AuthoredBy(okID)
	prefix := bot.Prefix(".")
	robot.On(prefix.Cmd("uptime").And(byID), cmdUptime(time.Now()))
	robot.On(prefix.Cmd("snap").And(byID), cmdHistSnapshot(collectChannel))
	robot.On(prefix.Cmd("invite").And(byID), cmdInvite)

	return nil
}
