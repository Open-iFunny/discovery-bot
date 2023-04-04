package main

import (
	"database/sql"
	"fmt"
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

					log.Trace("iter history")
					for event := range robot.Chat.IterMessages(compose.ListMessages(channel, 100, compose.NoPage())) {
						log.WithField("message_id", event.ID).Trace("enqueue event")
						event.Channel = channel
						events <- event
					}
				}
			}()
		}

		return nil
	}
}

const INSERT_CHUNK = 10_000

func snapEvents(events <-chan *ifunny.ChatEvent, procs int) func(*sql.DB, *bot.Bot) error {
	insertSnap := func(handle *sql.DB, buffer [][]any, errs chan error) {
		if err := insert(handle, "INSERT IGNORE INTO event_snap(id, event_type, channel, author, published) VALUES (?, ?, ?, ?, ?)", buffer); err != nil {
			errs <- err
		}
	}

	insertContent := func(handle *sql.DB, buffer [][]any, errs chan error) {
		if err := insert(handle, "INSERT IGNORE INTO event_message_content(id, content) VALUES (?, ?)", buffer); err != nil {
			errs <- err
		}
	}

	return func(handle *sql.DB, robots *bot.Bot) error {
		errs := make(chan error)

		for proc := 0; proc < procs; proc++ {
			go func() {
				for {
					bufferSnap := make([][]any, INSERT_CHUNK)
					bufferContent := make([][]any, INSERT_CHUNK)
					for index := range bufferSnap {
						switch event := <-events; true {
						case event == nil:
							errs <- fmt.Errorf("events stream closed")
						case event.Channel == "":
							errs <- fmt.Errorf("event has no channel: %+v", event)
						default:
							robots.Log.WithField("buffer_index", index).Trace("buffering event")
							bufferSnap[index] = []any{event.ID, event.Type, event.Channel, event.User.Nick, event.PubAt}
							bufferContent[index] = []any{event.ID, event.Text}
						}
					}

					robots.Log.Trace("writing buffers")
					go insertSnap(handle, bufferSnap, errs)
					go insertContent(handle, bufferContent, errs)
				}
			}()
		}

		return <-errs
	}
}
