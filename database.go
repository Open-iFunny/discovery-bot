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

	"github.com/go-sql-driver/mysql"
)

const (
	tMessageSnap      = "message_snap"
	tMessageSnapPlace = "message_snap_place"
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

type dblog struct{ logger *logrus.Logger }

func (log *dblog) Print(v ...interface{}) {
	log.logger.Trace(v...)
}

func dbSetup(logger *logrus.Logger) (*sql.DB, error) {
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

	mysql.SetLogger(&dblog{logger})
	handle.SetConnMaxLifetime(1 * time.Second)
	handle.SetMaxOpenConns(8)
	handle.SetMaxIdleConns(8)
	return handle, nil
}

func histSetPlace(handle *sql.DB, channel string, page int64, head string, finished bool) error {
	tx, err := handle.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	stmt, err := tx.Prepare("REPLACE INTO message_snap_place(channel, page, head, finished) VALUES (?, ?, ?, ?) LIMIT 1")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(channel, page, head, finished)
	return err
}

func histGetPlace(handle *sql.DB, channel string) (int64, string, bool, error) {
	page, head, finished := int64(0), "", false
	tx, err := handle.Begin()
	if err != nil {
		return page, head, finished, err
	}

	defer tx.Rollback()
	stmt, err := tx.Prepare("SELECT page, head, finished FROM message_snap_place WHERE channel = ? LIMIT 1")
	if err != nil {
		return page, head, finished, err
	}

	defer stmt.Close()
	result, err := stmt.Query(channel)
	if err != nil || !result.Next() {
		return page, head, finished, err
	}

	err = result.Scan(&page, &head, &finished)
	return page, head, finished, err
}

func histWrite(handle *sql.DB, head *ifunny.ChatEvent, buffer []*ifunny.ChatEvent, channel string) error {
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
		if event == nil {
			break
		}

		if _, err := stmt.Exec(event.ID, channel, event.User.Nick, int(event.PubAt), event.Text); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return histSetPlace(handle, channel, int64(head.PubAt), head.ID, false)
}

func doHistSeq(handle *sql.DB, channel string, robot *bot.Bot) error {
	log := robot.Log.WithFields(logrus.Fields{"channel": channel})
	pageIndex, _, finished, err := histGetPlace(handle, channel)
	if err != nil || finished {
		return err
	}

	page := compose.NoPage()
	if pageIndex != 0.0 {
		page = compose.Next(pageIndex)
	}

	index, bufLimit := 0, 2500
	buffer := make([]*ifunny.ChatEvent, bufLimit)
	iterEvent := robot.Chat.IterMessages(compose.ListMessages(channel, 100, page))
	for event := range iterEvent {
		if event == nil {
			log.Trace("iter end")
			goto flush
		}

		if event.Type != ifunny.TEXT_MESSAGE {
			continue
		}

		if index == bufLimit {
			log.Info("writing buffer")
			if err := histWrite(handle, event, buffer[:index], channel); err != nil {
				return err
			}

			index = 0
			buffer = make([]*ifunny.ChatEvent, bufLimit)
		}

		buffer[index] = event
		index++
	}

flush:
	if index != 0 {
		log.Info("flushing buffer")
		if err := histWrite(handle, buffer[index-1], buffer[:index], channel); err != nil {
			return err
		}
	}

	log.Info("marking complete")
	if err := histSetPlace(handle, channel, 0.0, "", true); err != nil {
		return err
	}

	return nil
}

func histSeq(rate time.Duration, channels <-chan string) func(robot *bot.Bot, handle *sql.DB) error {
	return func(robot *bot.Bot, handle *sql.DB) error {
		for channel := range channels {
			robot.Log.WithField("channel", channel).Info("hist seq GO")
			if err := doHistSeq(handle, channel, robot); err != nil {
				return err
			}

			robot.Log.WithField("channel", channel).Info("hist seq OK")
			<-time.After(rate)
		}

		return nil
	}
}
