package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/open-ifunny/discovery-bot/bot"
	"github.com/open-ifunny/discovery-bot/ifunny"
	"github.com/open-ifunny/discovery-bot/ifunny/compose"
)

func collectEventHist(rate time.Duration, channels <-chan string, events chan<- *ifunny.ChatEvent, procs int) func(*sql.DB, *bot.Bot) error {
	isFinished := func(handle *sql.DB, channel string) (bool, error) {
		finished := false
		err := query(handle, "SELECT finished FROM event_snap_place WHERE channel=?", []any{channel}, &finished)
		return finished, err
	}

	markFinished := func(handle *sql.DB, channel string) error {
		return query(handle, "REPLACE INTO event_snap_place(channel, finished) VALUES (?, ?)", []any{channel, true})
	}

	return func(handle *sql.DB, robot *bot.Bot) error {
		errs := make(chan error)

		for proc := 0; proc < procs; proc++ {
			go func() {
				for channel := range channels {
					log := robot.Log.WithField("channel", channel)
					if finished, err := isFinished(handle, channel); err != nil {
						errs <- err
					} else if finished {
						log.Info("already finished history iter")
						continue
					}

					log.Trace("iter history")
					for event := range robot.Chat.IterMessages(compose.ListMessages(channel, 500, compose.NoPage())) {
						log.WithField("message_id", event.ID).Trace("enqueue event")
						event.Channel = channel
						events <- event
					}

					if err := markFinished(handle, channel); err != nil {
						errs <- err
						return
					}
				}
			}()
		}

		return <-errs
	}
}

const INSERT_CHUNK = 100_000

func snapEvents(events <-chan *ifunny.ChatEvent, procs int) func(*sql.DB, *bot.Bot) error {
	insertSnap := func(handle *sql.DB, buffer [INSERT_CHUNK][]any) error {
		return insert(handle, "INSERT IGNORE INTO event_snap(id, event_type, channel, author, published) VALUES (?, ?, ?, ?, ?)", buffer)
	}

	insertContent := func(handle *sql.DB, buffer [INSERT_CHUNK][]any) error {
		return insert(handle, "INSERT IGNORE INTO event_message_content(id, content) VALUES (?, ?)", buffer)
	}

	return func(handle *sql.DB, robots *bot.Bot) error {
		errs := make(chan error)

		for proc := 0; proc < procs; proc++ {
			go func() {
				bufferSnap := [INSERT_CHUNK][]any{}
				bufferContent := [INSERT_CHUNK][]any{}

				for {
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
					if err := insertSnap(handle, bufferSnap); err != nil {
						errs <- err
					}

					if err := insertContent(handle, bufferContent); err != nil {
						errs <- err
					}
				}
			}()
		}

		return <-errs
	}
}
