package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-sql-driver/mysql"
)

const (
	tEventSnap       = "event_snap"
	tEventSnapPlace  = "event_snap_place"
	tMessageContent  = "event_message_content"
	tChannelSnap     = "channel_snap"
	tChannelSeqPlace = "channel_seq_place"
)

var tableDesc = [...][2]string{
	{
		tEventSnap,
		`id CHAR(32) UNIQUE PRIMARY KEY NOT NULL,
		event_type INT NOT NULL,
		channel CHAR(128) NOT NULL,
		author CHAR(32) NOT NULL,
		published BIGINT NOT NULL`,
	},
	{
		tEventSnapPlace,
		`channel CHAR(128) UNIQUE PRIMARY KEY NOT NULL,
		page BIGINT NOT NULL DEFAULT 0,
		head CHAR(32) NOT NULL DEFAULT '',
		finished BOOL NOT NULL DEFAULT FALSE`,
	},
	{
		tMessageContent,
		`id CHAR(32) UNIQUE PRIMARY KEY NOT NULL,
		content VARCHAR(1024) NOT NULL`,
	},
	{
		tChannelSnap,
		`name CHAR(255) UNIQUE PRIMARY KEY NOT NULL`,
	},
	{
		tChannelSeqPlace,
		`place CHAR(8) NOT NULL,
		finished BOOL NOT NULL DEFAULT FALSE,
		thread_lock INT UNIQUE PRIMARY KEY NOT NULL`,
	},
}
var connect = os.Getenv("IFUNNY_STATS_CONNECTION")

type dblog struct{ logger *logrus.Logger }

func (log *dblog) Print(v ...interface{}) {
	log.logger.Trace(v...)
}

func makeDB(logger *logrus.Logger) (*sql.DB, error) {
	handle, err := sql.Open("mysql", connect)
	if err != nil {
		return nil, err
	}

	if err = handle.Ping(); err != nil {
		return nil, err
	}

	for _, desc := range tableDesc {
		stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s) "+
			"ENGINE=MyISAM DEFAULT CHARSET='utf8' COLLATE='utf8_bin'",
			desc[0], desc[1])

		if _, err := handle.Query(stmt); err != nil {
			return nil, err
		}
	}

	threads, err := strconv.Atoi("IFUNNY_STATS_THREADS")
	if err != nil {
		threads = 32
	}

	mysql.SetLogger(&dblog{logger})
	handle.SetConnMaxLifetime(1 * time.Minute)
	handle.SetMaxOpenConns(threads)
	handle.SetMaxIdleConns(threads)
	return handle, nil
}

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
		if _, err := stmt.Exec(args...); err != nil {
			return err
		}

		return tx.Commit()
	}

	result, err := stmt.Query(args...)
	if err != nil || !result.Next() {
		return err
	}

	defer result.Close()
	return result.Scan(output...)
}

func insert(handle *sql.DB, query string, data [INSERT_CHUNK][]any) error {
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
	for _, dataRow := range data {
		if _, err := stmt.Exec(dataRow...); err != nil {
			return err
		}
	}

	return tx.Commit()
}
