package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/open-ifunny/discovery-bot/bot"
	"github.com/open-ifunny/discovery-bot/ifunny"
	"github.com/open-ifunny/discovery-bot/ifunny/compose"
)

var bearer = os.Getenv("IFUNNY_BEARER")
var userAgent = os.Getenv("IFUNNY_USER_AGENT")

func onChannelUpdate(_ *sql.DB, robot *bot.Bot) error {
	for {
		robot.Log.Trace("refresh channls subscribe")
		unsub, err := robot.Chat.OnChannelUpdate(func(eventType int, channel *ifunny.ChatChannel) error {
			switch eventType {
			case ifunny.EVENT_JOIN:
				robot.Subscribe(channel.Name)
				if channel.Type != compose.ChannelDM {
					return nil
				}

				return nil
			default:
				fmt.Printf("something else happened [%d]: %+v", eventType, channel)
			}

			return nil
		})

		if err != nil {
			return err
		}

		<-time.Tick(15 * time.Minute)
		unsub()
	}
}

func onChannelInvite(_ *sql.DB, robot *bot.Bot) error {
	_, err := robot.Chat.OnChannelInvite(func(eventType int, channel *ifunny.ChatChannel) error {
		return robot.Chat.Call(compose.InviteResponse(channel.Name, true), nil)
	})

	return err
}

var threadChannelSeq = func() int {
	if v, err := strconv.Atoi(os.Getenv("IFUNNY_CHANNEL_SEQ_THREADS")); err != nil {
		return 8
	} else {
		return v
	}
}()

var threadEventHist = func() int {
	if v, err := strconv.Atoi(os.Getenv("IFUNNY_EVENT_HIST_THREADS")); err != nil {
		return 64
	} else {
		return v
	}
}()

var threadEventSnap = func() int {
	if v, err := strconv.Atoi(os.Getenv("IFUNNY_EVENT_SNAP_THREADS")); err != nil {
		return 8
	} else {
		return v
	}
}()

var collectChannel = make(chan string, 32)
var collectEvent = make(chan *ifunny.ChatEvent)
var forevers = [...]struct {
	name string
	call func(*sql.DB, *bot.Bot) error
}{
	{"onCommand", onCommand},
	{"onChannelUpdate", onChannelUpdate},
	{"onChannelInvite", onChannelInvite},
	{"collectChannelSeq", collectChannelSeq(10*time.Millisecond, collectChannel, threadChannelSeq, 0)},
	{"collectEventHist", collectEventHist(10*time.Millisecond, collectChannel, collectEvent, threadEventHist)},
	{"snapEvents", snapEvents(collectEvent, threadEventSnap)},
}

var tickers = [...]struct {
	name     string
	interval time.Duration
	tick     func(*sql.DB, *bot.Bot) error
}{
	{"collect-channel-trending", 1 * time.Hour, collectChannelTrending(100*time.Millisecond, collectChannel)},
}

func main() {
	robot, err := bot.MakeBot(bearer, userAgent)
	if err != nil {
		panic("error in makeBot: " + err.Error())
	}

	handle, err := makeDB(robot.Log)
	if err != nil {
		panic("error in makeDB: " + err.Error())
	}

	for _, forever := range forevers {
		go func(name string, f func(*sql.DB, *bot.Bot) error) {
			robot.Log.Infof("call forever %s", name)

			if err := f(handle, robot); err != nil {
				panic(fmt.Sprintf("error in forever: %s: %s", name, err))
			}
		}(forever.name, forever.call)
	}

	for _, ticker := range tickers {
		go func(name string, interval time.Duration, call func(*sql.DB, *bot.Bot) error) {
			for iter := 0; ; iter++ {
				robot.Log.WithField("iter", iter).Infof("call ticker %s", name)

				if err := call(handle, robot); err != nil {
					panic(fmt.Sprintf("error in ticker[iter: %d]: %s: %s", iter, name, err))
				}

				<-time.Tick(interval)
			}
		}(ticker.name, ticker.interval, ticker.tick)
	}

	for {
		robot.Log.Info("listening")
		if err := robot.Listen(); err != nil {
			robot.Log.Errorf("error in listen: %s", err)
		}

		<-time.After(5 * time.Second)
		robot.Log.Info("reconnecting")
		chat, err := robot.Client.Chat()
		if err != nil {
			robot.Log.Errorf("err reconnecting: %s", err)
		}

		robot.Chat = chat
	}
}
