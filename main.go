package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gastrodon/popplio/bot"
	"github.com/gastrodon/popplio/ifunny"
	"github.com/gastrodon/popplio/ifunny/compose"
)

var bearer = os.Getenv("IFUNNY_BEARER")
var userAgent = os.Getenv("IFUNNY_USER_AGENT")

func joinTrending(robot *bot.Bot) error {
	channels, err := robot.Client.GetChannels(compose.ChatsTrending)
	if err != nil {
		return err
	}

	for _, channel := range channels {
		if channel.MembersTotal == 500 || channel.JoinState == compose.Joined {
			fmt.Printf("skipping %s\n", channel.Name)
			continue
		}

		if err := robot.Chat.Call(compose.JoinChannel(channel.Name), nil); err != nil {
			return err
		}
	}

	return nil
}

func onChannelUpdate(robot *bot.Bot) error {
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

		<-time.Tick(15 * time.Second)
		unsub()
	}
}

func onChannelInvite(robot *bot.Bot) error {
	_, err := robot.Chat.OnChannelInvite(func(eventType int, channel *ifunny.ChatChannel) error {
		return robot.Chat.Call(compose.Invite(channel.Name, true), nil)
	})

	return err
}

var histChan = make(chan string)
var registers = [...]func(*bot.Bot) error{
	onChannelUpdate, onChannelInvite,
	collectSeq(25*time.Millisecond, histChan),
	onCommand,
}

var tickers = [...]struct {
	interval time.Duration
	tick     func(*bot.Bot) error
}{
	{1 * time.Hour, joinTrending},
}

func main() {
	robot, err := bot.MakeBot(bearer, userAgent)
	if err != nil {
		panic(err)
	}

	for _, reg := range registers {
		go func(reg func(*bot.Bot) error) {
			if err := reg(robot); err != nil {
				robot.Log.Error("error: " + err.Error())
			}
		}(reg)
	}

	for _, ticker := range tickers {
		go func(tRobot *bot.Bot, interval time.Duration, call func(*bot.Bot) error) {
			for {
				if err := call(tRobot); err != nil {
					panic(err)
				}

				<-time.Tick(interval)
			}
		}(robot, ticker.interval, ticker.tick)
	}

	handle, err := dbSetup(robot.Log)
	if err != nil {
		panic(err)
	}

	for i := 0; i != 15; i++ {
		go func() {
			if err := histSeq(1500*time.Millisecond, histChan)(robot, handle); err != nil {
				panic(err)
			}
			<-time.After(100 * time.Millisecond)
		}()
	}

	if err := robot.Listen(); err != nil {
		panic(err)
	}
}
