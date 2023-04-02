package main

import (
	"time"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny/compose"
)

var channelRune = []rune("abcdefghijklmnopqrstuvwxyz1234567890_")

func collectSeq(rate time.Duration, channels chan<- string) func(robot *bot.Bot) error {
	return func(robot *bot.Bot) error {
		for _, first := range channelRune {
			for _, second := range channelRune {
				for _, third := range channelRune {
					query := string([]rune{first, second, third})
					log := robot.Log.WithField("query", query)
					log.Trace("iter results")

					for channel := range robot.Client.IterChannels(compose.ChatsQuery(query, 100, compose.SPage{})) {
						log.WithField("channel", channel.Name).Info("enqueue channel result")
						channels <- channel.Name
					}

					<-time.Tick(rate)
				}
			}
		}

		robot.Log.Info("finished enqueuing every seq query OK")
		return nil
	}
}
