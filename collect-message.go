package main

import (
	"database/sql"
	"time"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny"
	"github.com/gastrodon/popplio/ifunny/compose"
)

func collectEventHist(rate time.Duration, channels <-chan string, events chan<- *ifunny.ChatEvent, procs int) func(*sql.DB, *bot.Bot) error {
	return func(_ *sql.DB, robot *bot.Bot) error {
		for proc := 0; proc < procs; proc++ {
			go func() {
				for channel := range channels {
					log := robot.Log.WithField("channel", channel)

					log.Info("iter history")
					for event := range robot.Chat.IterMessages(compose.ListMessages(channel, 100, compose.NoPage())) {
						log.WithField("message_id", event.ID).Trace("enqueue event")
						events <- event
					}
				}
			}()
		}

		return nil
	}
}
