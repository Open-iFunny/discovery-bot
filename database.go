package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny"
	"github.com/gastrodon/popplio/ifunny/compose"
	"github.com/sirupsen/logrus"
)

const (
	tMessageSnap      = "message_snap"
	tMessageSnapPlace = "message_snap_place"

	nothing = 0 * time.Nanosecond
)

var tableDesc = [...][2]string{
	{
		tMessageSnap,
		`id CHAR(32) UNIQUE PRIMARY KEY NOT NULL,
		channel CHAR(128) NOT NULL,
		author CHAR(128) NOT NULL,
		published BIGINT NOT NULL,
		content VARCHAR(4096) NOT NULL`,
	},
	{
		tMessageSnapPlace,
		`channel CHAR(128) UNIQUE PRIMARY KEY NOT NULL,
		page BIGINT NOT NULL,
		head CHAR(32) NOT NULL,
		finished BOOL NOT NULL DEFAULT FALSE`,
	},
}
var connect = os.Getenv("IFUNNY_STATS_CONNECTION")

func dbSetup() (*sql.DB, error) {
	handle, err := sql.Open("mysql", connect)
	if err != nil {
		return nil, err
	}

	if err = handle.Ping(); err != nil {
		return nil, err
	}

	for _, desc := range tableDesc {
		if _, err := handle.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", desc[0], desc[1])); err != nil {
			return nil, err
		}
	}

	return handle, nil
}

func histSetPlace(handle *sql.DB, channel string, page int64, head string, finished bool) error {
	_, err := handle.Query("REPLACE INTO message_snap_place(channel, page, head, finished) VALUES (?, ?, ?, ?)", channel, page, head, finished)
	return err
}

func histGetPlace(handle *sql.DB, channel string) (int64, string, bool, error) {
	result, err := handle.Query("SELECT page, head, finished FROM message_snap_place WHERE channel = ?", channel)
	if err != nil {
		return 0, "", false, err
	}

	if !result.Next() {
		fmt.Println("no results, not done")
		return 0, "", false, nil
	}

	page, head, finished := int64(0), "", false
	err = result.Scan(&page, &head, &finished)

	fmt.Println("query results", page, head, finished, err)
	return page, head, finished, err
}

func histWrite(handle *sql.DB, buffer []*ifunny.ChatEvent, channel string) error {
	tx, err := handle.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT IGNORE INTO message_snap(id, channel, author, published, content) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	for _, event := range buffer {
		if _, err := stmt.Exec(event.ID, channel, event.User.Nick, int(event.PubAt), event.Text); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func histSeq(rate time.Duration, channels <-chan string) func(handle *sql.DB, robot *bot.Bot) error {
	return func(handle *sql.DB, robot *bot.Bot) error {
		for name := range channels {
			log := robot.Log.WithFields(logrus.Fields{"channel": name})
			log.Info("hist seq GO")

			pageIndex, _, finished, err := histGetPlace(handle, name)
			if err != nil {
				log.Error(err)
				return err
			}

			if finished {
				log.Info("hist seq already done")
				continue
			}

			page := compose.NoPage()
			if pageIndex != 0.0 {
				page = compose.Next(pageIndex)
			}

			index, bufLimit := 0, 10_000
			buffer := make([]*ifunny.ChatEvent, bufLimit)
			for event := range robot.Chat.IterMessages(compose.ListMessages(name, 100, page)) {
				if index == bufLimit {
					log.Trace("writing buffer")

					if err := histWrite(handle, buffer, name); err != nil {
						log.Error("err writing buffer: " + err.Error())
						return err
					}

					if err := histSetPlace(handle, name, int64(event.PubAt), event.ID, false); err != nil {
						log.Error("err writing place: " + err.Error())
						return err
					}

					index = 0
					buffer = make([]*ifunny.ChatEvent, bufLimit)
				}

				buffer[index] = event
				index++
				if rate != nothing {
					<-time.Tick(rate)
				}
			}

			if index != 0 {
				log.Trace("writing final buffer")

				if err := histWrite(handle, buffer[:index], name); err != nil {
					log.Error("err writing buffer: " + err.Error())
					return err
				}

				if err := histSetPlace(handle, name, 0.0, "", true); err != nil {
					log.Error("err writing place: " + err.Error())
					return err
				}
			}

			log.Info("hist seq OK")
		}

		return nil
	}
}
