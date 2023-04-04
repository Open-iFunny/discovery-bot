package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-sql-driver/mysql"
)

const (
	tMessageSnap      = "message_snap"
	tMessageContent   = "message_content"
	tMessageSnapPlace = "message_snap_place"
	tChannelSnap      = "channel_snap"
	tChannelSeqPlace  = "channel_seq_place"
)

var tableDesc = [...][2]string{
	{
		tMessageSnap,
		`id CHAR(32) UNIQUE PRIMARY KEY NOT NULL,
		channel CHAR(128) NOT NULL,
		author CHAR(32) NOT NULL,
		published BIGINT NOT NULL`,
	},
	{
		tMessageContent,
		`id CHAR(32) UNIQUE PRIMARY KEY NOT NULL,
		content VARCHAR(1024) NOT NULL`,
	},
	{
		tMessageSnapPlace,
		`channel CHAR(128) UNIQUE PRIMARY KEY NOT NULL,
		page BIGINT NOT NULL,
		head CHAR(32) NOT NULL,
		finished BOOL NOT NULL DEFAULT FALSE`,
	},
	{
		tChannelSnap,
		`name CHAR(255) UNIQUE PRIMARY KEY NOT NULL`,
	},
	{
		tChannelSeqPlace,
		`place CHAR(8) NOT NULL,
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

	mysql.SetLogger(&dblog{logger})
	handle.SetConnMaxLifetime(1 * time.Second)
	handle.SetMaxOpenConns(8)
	handle.SetMaxIdleConns(8)
	return handle, nil
}
