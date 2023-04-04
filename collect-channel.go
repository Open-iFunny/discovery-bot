package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny/compose"
	"github.com/sirupsen/logrus"
)

const PLACE_INDEX = 1

func query(handle *sql.DB, query string, args []any, output ...interface{}) error {
	tx, err := handle.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	if len(output) == 0 {
		_, err := stmt.Exec(args...)
		return err
	}

	result, err := stmt.Query(args...)
	if err != nil || !result.Next() {
		return err
	}

	defer result.Close()
	return result.Scan(output...)
}

var channelRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890_")

func countPlacePart(part rune) (int, error) {
	for index, each := range channelRunes {
		if each == part {
			return index, nil
		}
	}

	return 0, fmt.Errorf("rune %s not in placeable channelRunes", string(part))
}

func countPlace(place string) ([3]int, error) {
	parts := [3]int{}
	for index, each := range place {
		part, err := countPlacePart(each)
		if err != nil {
			return parts, err
		}

		parts[index] = part
	}

	return parts, nil
}

func collectChannelSeq(rate time.Duration, channels chan<- string, lock int) func(*sql.DB, *bot.Bot) error {
	getPlace := func(handle *sql.DB) ([3]int, error) {
		place := ""
		if err := query(handle, "SELECT place FROM channel_seq_place WHERE thread_lock=?", []any{lock}, &place); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}

			return [...]int{0, 0, 0}, err
		}

		return countPlace(place)
	}

	setPlace := func(handle *sql.DB, place string) error {
		return query(handle, "REPLACE INTO channel_seq_place(place, thread_lock) VALUES (?, ?)", []any{place, lock})
	}

	iterQuery := func(place [3]int) <-chan string {
		result := make(chan string)

		go func() {
			for _, first := range channelRunes[place[0]:] {
				for _, second := range channelRunes[place[1]:] {
					for _, third := range channelRunes[place[2]:] {
						result <- string([]rune{first, second, third})
					}
				}
			}

			close(result)
		}()

		return result
	}

	return func(handle *sql.DB, robot *bot.Bot) error {
		place, err := getPlace(handle)
		if err != nil {
			return err
		}

		for query := range iterQuery(place) {
			log := robot.Log.WithFields(logrus.Fields{"start_place": place, "query": query})

			log.Trace("set place")
			if err := setPlace(handle, query); err != nil {
				return err
			}

			log.Info("iter results")
			for channel := range robot.Client.IterChannels(compose.ChatsQuery(query, 100, compose.SPage{})) {
				log.WithField("channel", channel.Name).Trace("enqueue channel result")
				channels <- channel.Name
				<-time.Tick(rate)
			}
		}

		return nil
	}
}

func collectChannelTrending(rate time.Duration, channels chan<- string) func(*sql.DB, *bot.Bot) error {
	return func(_ *sql.DB, robot *bot.Bot) error {
		if trending, err := robot.Client.GetChannels(compose.ChatsTrending); err != nil {
			return err
		} else {
			for _, channel := range trending {
				channels <- channel.Name
				<-time.Tick(rate)
			}
		}

		return nil
	}
}
